package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
	"fmt"
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
	log.Println(versionStr)
	if err := fortiC.Init(); err != nil {
		log.Fatalln(err)
	}
	// new subprocess
	vpnProg := exec.Command(fortiC.BinaryPath, "-s", fortiC.ServerAddr, "-u", fortiC.Username, "-p")
	vpnStdErr, err := vpnProg.StderrPipe()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("stderr pipe got.")
	vpnStdOut, err := vpnProg.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("stdout pipe got.")
	vpnStdIn, err := vpnProg.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("stdin pipe got.")
	// start new process
	err = vpnProg.Start()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("vpn process started.")
	// start input and output data
	// stdout, stdin
	go func() {
		scnr := bufio.NewScanner(vpnStdOut)
		scnr.Split(bufio.ScanLines)
		for scnr.Scan() {
			curLine := scnr.Bytes()
			// never close fxxking stdin!
			if bytes.HasPrefix(curLine, []byte("password:")) {
				_, _ = vpnStdIn.Write([]byte(fortiC.Password + "\n"))
				log.Println("write password to vpn cli stdin.")
				fmt.Println(string(curLine))
				continue
			}
			if bytes.Contains(curLine, []byte("Confirm (y/n) [default=n]:")) {
				_, _ = vpnStdIn.Write([]byte(fortiC.insecureAns + "\n"))
				log.Printf("answered %s to cert insecure warning. \n", fortiC.insecureAns)
				fmt.Println(string(curLine))
				continue
			}
			fmt.Println(string(curLine))
		}
	}()
	// stderr
	go func() {
		_, err = io.Copy(os.Stderr, vpnStdErr)
		if errors.Is(err, io.EOF) {
			return
		} else {
			log.Println("Error for Stdout Redirect:", err)
		}
	}()
	// handle signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGKILL)
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
		log.Println(err)
	}
	log.Println("vpn process killed. exit.")
}
