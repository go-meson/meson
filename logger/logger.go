package logger

import (
	"bufio"
	"log"
	"os"
	"syscall"
)

func SetFileLogger(name string) error {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	log.SetOutput(f)
	return nil
}

func redirectIO(f *os.File, prefix string) {
	go func() {
		f := f
		pr, pw, err := os.Pipe()
		if err != nil {
			log.Fatal(err)
			return
		}
		defer pr.Close()
		defer pw.Close()
		err = syscall.Dup2(int(pw.Fd()), int(f.Fd()))
		if err != nil {
			log.Fatal(err)
			return
		}
		r := bufio.NewReader(pr)
		for {
			s, err := r.ReadString('\n')
			if err != nil {
				break
			}
			log.Printf("*%s*:%s", prefix, s)
		}

	}()
}

func RedirectStdout() {
	redirectIO(os.Stdout, "STDOUT")
}

func RedirectStderr() {
	redirectIO(os.Stderr, "STDERR")
}
