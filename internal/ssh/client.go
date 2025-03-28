package ssh

import (
	"fmt"
	"net"

	"github.com/Michaelpalacce/go-btva/pkg/prompt"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

// GetClient based on the application's options will create a new client with ssh to the given machine
func GetClient(vmIP, username, password, privateKey, privateKeyPassphrase string) (*goph.Client, error) {
	var (
		client *goph.Client
		auth   goph.Auth
		err    error
	)

	if privateKey != "" {
		// Start new ssh connection with private key.
		auth, err = goph.Key(privateKey, privateKeyPassphrase)
		if err != nil {
			return nil, fmt.Errorf("there was an error using the private key %s, err was: %w", privateKey, err)
		}
	} else {
		auth = goph.Password(password)
	}

	if client, err = goph.NewConn(&goph.Config{User: username, Addr: vmIP, Port: 22, Auth: auth, Timeout: goph.DefaultTimeout, Callback: VerifyHost}); err != nil {
		return nil, fmt.Errorf("there was an error while establishing connection to remote server for infra setup. Err was: %w", err)
	}

	return client, nil
}

// VerifyHost will check if the host is trusted or the signature has changed. If neither, prompts the user to accept the connection
func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {
	hostFound, err := goph.CheckKnownHost(host, remote, key, "")

	// Host in known hosts but key mismatch!
	// Maybe because of MAN IN THE MIDDLE ATTACK!
	if hostFound && err != nil {
		return err
	}

	if hostFound && err == nil {
		return nil
	}

	if host == "" {
		return fmt.Errorf("host is an empty string.")
	}

	if askIsHostTrusted(host, key) == false {
		return fmt.Errorf("you didn't trust host %s, exiting.", host)
	}

	return goph.AddKnownHost(host, remote, key, "")
}

func askIsHostTrusted(host string, key ssh.PublicKey) bool {
	a, err := prompt.AskText(fmt.Sprintf("Unknown Host: %s \nFingerprint: %s \n", host, ssh.FingerprintSHA256(key)), fmt.Sprintf("Would you like to add it? type yes or no: "))
	if err != nil {
		return false
	}

	return prompt.IsYesAnswer(a)
}
