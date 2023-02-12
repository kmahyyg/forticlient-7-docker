//go:build unix

package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"github.com/creack/pty"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var (
	fortiC = &fortiConfig{}

	//go:embed version.txt
	versionStr string
)

var (
	ErrNotPodman = errors.New("this program needs to be run inside container managed by podman, " +
		"docker or others are not supported due to limitation by Forticlient ")
	ErrRequirementNotSatisfied = errors.New("please set all required envvar, read README")
)

type fortiConfig struct {
	Username      string
	Password      string
	AllowInsecure bool
	insecureAns   string
	ServerAddr    string
	BinaryPath    string
	// reserved, not in-use
	TOTP        string
	CertPath    string
	CertKeyPath string
}

func (fc *fortiConfig) Init() error {
	fc.BinaryPath = os.Getenv("FORTIVPN_CLI")
	if fc.BinaryPath != "" {
		_, err := os.Stat(fc.BinaryPath)
		if err != nil {
			log.Println("vpncli may not exist.")
			return ErrRequirementNotSatisfied
		}
		log.Println("vpncli found.")
	} else {
		log.Println("vpncli cannot be found from envvar")
		return ErrRequirementNotSatisfied
	}
	fc.Username = os.Getenv("FORTIVPN_USR")
	fc.ServerAddr = os.Getenv("FORTIVPN_SRV")
	fc.Password = os.Getenv("FORTIVPN_PASSWD")
	fc.AllowInsecure = func() bool {
		if os.Getenv("ALLOW_INSECURE") != "" {
			fc.insecureAns = "y"
			return true
		}
		fc.insecureAns = "n"
		return false
	}()
	if fc.Username == "" || fc.ServerAddr == "" || fc.Password == "" {
		log.Println("fc init failed, any of requirement is empty.")
		return ErrRequirementNotSatisfied
	}
	log.Println("fc init done.")
	return nil
}

func init() {
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lmicroseconds | log.Lshortfile)
	// thanks to https://systemd.io/CONTAINER_INTERFACE/
	if os.Getenv("container") != "podman" {
		log.Fatalln(ErrNotPodman)
	}
	log.Println("running in podman detected.")
}

func main() {
	// start
	log.Println("version: ", versionStr)
	if err := fortiC.Init(); err != nil {
		log.Fatalln(err)
	}
	var err error
	// Debug:
	log.Printf("config: %+v \n", *fortiC)
	// new subprocess
	vpnProg := exec.Command(fortiC.BinaryPath, "-s", fortiC.ServerAddr, "-u", fortiC.Username, "-p")
	vpnProg.SysProcAttr = &syscall.SysProcAttr{
		Setsid:  true,
		Setctty: true,
	}
	ptyInst, ttyInst, err := pty.Open()
	if err != nil {
		panic(err)
	}
	defer ttyInst.Close()
	defer ptyInst.Close()
	// start input and output data
	vpnProg.Stdin = ttyInst
	vpnProg.Stdout = ttyInst
	vpnProg.Stderr = ttyInst
	// stdout, stdin
	go func() {
		scnr := bufio.NewScanner(ptyInst)
		scnr.Split(bufio.ScanWords)
		for scnr.Scan() {
			curLine := scnr.Bytes()
			log.Println("scanned output from stdout: ", string(curLine))
			// never close fxxking stdin!
			if bytes.HasPrefix(curLine, []byte("password:")) {
				_, _ = io.WriteString(ptyInst, fortiC.Password+"\n")
				log.Println("write password to vpn cli stdin.")
				continue
			}
			if bytes.Contains(curLine, []byte("[default=n]:Confirm")) {
				_, _ = io.WriteString(ptyInst, fortiC.insecureAns+"\n")
				log.Printf("answered %s to cert insecure warning. \n", fortiC.insecureAns)
				break
			}
		}
		log.Println("scan input completed, all answer finished. Now direct copy stdout and show.")
		_, err = io.Copy(os.Stdout, ptyInst)
		if errors.Is(err, io.EOF) {
			return
		} else {
			log.Println("Error for Stdout Redirect:", err)
		}
	}()
	// start new process
	err = vpnProg.Start()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("vpn process started.")
	// handle signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	// handle exit
	// wait until close
	go func() {
		err := vpnProg.Wait()
		if err != nil {
			log.Println("vpnprog unexpected exit, err: ", err)
		}
		log.Println("vpnprog exit unexpectedly.")
		sigChan <- syscall.SIGTERM
	}()
	<-sigChan
	log.Println("received signal of killing vpn process, now clean up.")
	err = vpnProg.Process.Kill()
	if err != nil {
		log.Println("process kill err: ", err)
	}
	log.Println("vpn process killed. exit.")
}
