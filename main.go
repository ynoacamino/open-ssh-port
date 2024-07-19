package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
)

const (
	PASSWORD = "user10"
)

func main() {
	openSSH()

	username := getUser()

	setPassword(username, PASSWORD)
	getIp()
}

func getIp() {
	// Obtener las interfaces de red disponibles
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error al obtener interfaces:", err)
		return
	}

	// Iterar sobre las interfaces de red
	for _, i := range interfaces {
		// Ignorar interfaces que no están activas o no son de tipo hardware
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Obtener las direcciones asociadas a la interfaz
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println("Error al obtener direcciones para la interfaz:", err)
			continue
		}

		// Iterar sobre las direcciones de la interfaz
		for _, addr := range addrs {
			// Convertir la dirección a formato IPNet
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Ignorar direcciones IPv6 y direcciones de loopback
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			// Imprimir la dirección IP de la red privada
			fmt.Println("-----------------------------------------------")
			fmt.Println("Dirección IP de la red privada:", ip.String())
			fmt.Println("-----------------------------------------------")
		}
	}
}

func openSSH() {
	// El comando que quieres ejecutar
	cmd := exec.Command("cmd.exe", "/c", "Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0")

	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("La salida del comando fue: %s\n", output)

	cmd = exec.Command("cmd.exe", "/c", "Set-Service -Name sshd -StartupType 'Automatic'")
	output, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("La salida del comando fue: %s\n", output)

	cmd = exec.Command("cmd.exe", "/c", "Start-Service sshd")
	output, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("La salida del comando fue: %s\n", output)

	cmd = exec.Command("cmd.exe", "/c", "New-NetFirewallRule -Name sshd -DisplayName 'OpenSSH Server (sshd)' -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort 22")
	output, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("La salida del comando fue: %s\n", output)

	cmd = exec.Command("cmd.exe", "/c", "Get-Service -Name sshd")
	output, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("La salida del comando fue: %s\n", output)

	cmd = exec.Command("cmd.exe", "/c", "Test-NetConnection -ComputerName localhost -Port 22")
	output, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("La salida del comando fue: %s\n", output)
}

func getUser() string {
	// Obtener el nombre de usuario actual
	username, err := exec.Command("cmd.exe", "/c", "echo %USERNAME%").Output()
	if err != nil {
		fmt.Println("Error al obtener nombre de usuario:", err)
		log.Fatal(err)
	}

	// Imprimir el nombre de usuario
	fmt.Println("Nombre de usuario:", string(username))

	return string(username)
}

func setPassword(username string, pswd string) {
	// Cambiar la contraseña de un usuario
	cmd := exec.Command("cmd.exe", "/c", "net user "+username+" "+pswd)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error al cambiar contraseña:", err)
		log.Fatal(err)
	}

	// Imprimir la salida del comando
	fmt.Printf("La salida del comando fue: %s\n", output)

	fmt.Println("-----------------------------------------------")
	fmt.Println("Contraseña cambiada con éxito")
	fmt.Println("Nombre de usuario:", username)
	fmt.Println("Contraseña:", pswd)
	fmt.Println("-----------------------------------------------")
}
