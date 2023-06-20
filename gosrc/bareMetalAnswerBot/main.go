package main

import (
	_ "embed"
	"errors"
	expect "github.com/Netflix/go-expect"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
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
}

func main() {
	// start
	log.Println("version: ", versionStr)
	if err := fortiC.Init(); err != nil {
		log.Fatalln(err)
	}

	// Debug:
	log.Printf("config: %+v \n", *fortiC)

	// pre-flight check and prepare
	consoleProg, err := expect.NewConsole(expect.WithStdout(os.Stdout), expect.WithStdin(os.Stdin))
	if err != nil {
		log.Fatalln(err)
	}
	defer consoleProg.Close()

	// new subprocess
	vpnCmd := exec.Command(fortiC.BinaryPath, "-u", fortiC.Username, "-s", fortiC.ServerAddr, "-p")
	vpnCmd.Stdin = consoleProg.Tty()
	vpnCmd.Stderr = consoleProg.Tty()
	vpnCmd.Stdout = consoleProg.Tty()

	go func() {
		expectPwd := expect.String("password:")
		expectCert := expect.String("Confirm (y/n) [default=n]:")
		for {
			data, err := consoleProg.Expect(expectPwd, expectCert)
			if err != nil {
				log.Println("expect err: ", err)
				break
			}
			if strings.Contains(data, "password:") {
				_, _ = consoleProg.SendLine(fortiC.Password)
				continue
			} else if strings.Contains(data, "Confirm (y/n) [default=n]:") {
				_, _ = consoleProg.SendLine(fortiC.insecureAns)
				break
			} else {
				continue
			}
		}
		consoleProg.ExpectEOF()
	}()

	err = vpnCmd.Start()
	if err != nil {
		log.Fatalln(err)
	}
	// handle signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	// handle exit
	// wait until close
	go func() {
		err := vpnCmd.Wait()
		if err != nil {
			log.Println("vpnprog exit unexpectedly: ", err)
		} else {
			// unpredicted exit, but wait returns nothing.
			log.Println("vpnprog got signal to exit, now exited.")
			sigChan <- syscall.SIGABRT
		}
	}()
	<-sigChan
	log.Println("received signal of killing vpn process, now clean up.")
	err = vpnCmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		log.Println(err)
	}
	log.Println("sleep 5 seconds for cleanup.")
	time.Sleep(5 * time.Second)
	log.Println("vpn process killed. exit.")
}
