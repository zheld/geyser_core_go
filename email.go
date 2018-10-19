package core

import (
    "crypto/tls"
    "fmt"
    "net/smtp"
    "strings"
)

type mail struct {
    senderId string
    toIds    []string
    subject  string
    body     string
}

type smtpServer struct {
    host string
    port string
}

func (s *smtpServer) serverMailName() string {
    return s.host + ":" + s.port
}

func (mail *mail) buildMailMessage() string {
    message := ""
    message += fmt.Sprintf("From: %s\r\n", mail.senderId)
    if len(mail.toIds) > 0 {
        message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
    }

    message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
    message += "\r\n" + mail.body

    return message
}

func SendMail(addr, pswd, subject, msg, recipient string) error {
    mail := mail{}
    mail.senderId = addr
    mail.toIds = []string{recipient}
    mail.subject = subject
    mail.body = msg

    messageBody := mail.buildMailMessage()

    smtpServer := smtpServer{host: "smtp.gmail.com", port: "465"}

    //build an auth
    auth := smtp.PlainAuth("", mail.senderId, pswd, smtpServer.host)

    // Gmail will reject connection if it's not secure
    // TLS config
    tlsconfig := &tls.Config{
        InsecureSkipVerify: true,
        ServerName:         smtpServer.host,
    }

    conn, err := tls.Dial("tcp", smtpServer.serverMailName(), tlsconfig)
    if err != nil {
        return err
    }

    client, err := smtp.NewClient(conn, smtpServer.host)
    if err != nil {
        return err
    }

    // step 1: Use Auth
    if err = client.Auth(auth); err != nil {
        return err
    }

    // step 2: add all from and to
    if err = client.Mail(mail.senderId); err != nil {
        return err
    }
    for _, k := range mail.toIds {
        if err = client.Rcpt(k); err != nil {
            return err
        }
    }

    // Data
    w, err := client.Data()
    if err != nil {
        return err
    }

    _, err = w.Write([]byte(messageBody))
    if err != nil {
        return err
    }

    err = w.Close()
    if err != nil {
        return err
    }

    client.Quit()

    return nil
}
