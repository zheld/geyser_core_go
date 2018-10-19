package core

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "io/ioutil"
    "fmt"
    "encoding/hex"
    "crypto/aes"
    "io"
    "crypto/cipher"
)

var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAwOJK1RJBUwRu/5aCyktTaietXFMOAAkElhSq1M6BocVWs7yD
y592CX30Bl0Ul4faWM9EZSlhak8Ay1CdMNis+lBZanKmAO2bPmSIIYBDQE2BzLIo
MoJWi/Cd5PevioKSRPytqVB/S4+xz1IOD8Y407SZM3LfZ5XMfqC+VHpcnAycQ8iT
FK0s3yjImathFNF3U7fiEzU4G7PJRn8e9ntubDd1pXYABqrVF/REcd/3Rs/qrlhG
v3b7tAXZb2lkSLdCq3Md+BMksxUCoH3rZijSphbZSCdIrzofg+IG0y5WtdsBz6uw
Ol2QX/EUoEdO+xhLgaOFykUoWz037ZzkLEhKkQIDAQABAoIBAB+1lAPPSnnxYqYW
Ak5rb70l5LQm20haMyzRHPx7Loh/vq8xsKELCAardDCPoNEAfn7XJDFVSjSF5GWI
TS84j8de6jQ7wNqqNTleoZqQUX4Cv/H83+rdzoiW9/4qUet9Z7p7p7kMCMFNUDf7
D2C8f58eM4lnux52W/X9SwzsSMlGaGHcAKPz4vXUFWyt3naVtANhdkHjgKxA0Ev4
W7yRgpbOKruPKzBNTRXAqq+yHZj/pONtXl8do+plwhHU8CW0BPyvkU4DH52lxWza
mM71ow8UJC30FXF/NZ+wthFnRZO3/dhaeuNYgX7yAs3DhNn7Q8nzU4ujd8ug2OGf
iJ9C8YECgYEA32KthV7VTQRq3VuXMoVrYjjGf4+z6BVNpTsJAa4kF+vtTXTLgb4i
jpUrq6zPWZkQ/nR7+CuEQRUKbky4SSHTnrQ4yIWZTCPDAveXbLwzvNA1xD4w4nOc
JgG/WYiDtAf05TwC8p/BslX20Ox8ZAXUq6pkAeb1t8M2s7uDpZNtBMkCgYEA3QuU
vrrqYuD9bQGl20qJI6svH875uDLYFcUEu/vA/7gDycVRChrmVe4wU5HFErLNFkHi
uifiHo75mgBzwYKyiLgO5ik8E5BJCgEyA9SfEgRHjozIpnHfGbTtpfh4MQf2hFsy
DJbeeRFzQs4EW2gS964FK53zsEtnr7bphtvfY4kCgYEAgf6wr95iDnG1pp94O2Q8
+2nCydTcgwBysObL9Phb9LfM3rhK/XOiNItGYJ8uAxv6MbmjsuXQDvepnEp1K8nN
lpuWN8rXTOG6yG1A53wWN5iK0WrHk+BnTA7URcwVqJzAvO3RYVPqqlcwTKByOtrR
yhxcGmdHMusdWDaVA7PpS1ECgYATCGs/XQLPjsXje+/XCPzz+Epvd7fi12XpwfQd
Z5j/q82PsxC+SQCqR38bwwJwELs9/mBSXRrIPNFbJEzTTbinswl9YfGNUbAoT2AK
GmWz/HBY4uBoDIgEQ6Lu1o0q05+zV9LgaKExVYJSL0EKydRQRUimr8wK0wNTivFi
rk322QKBgHD3aEN39rlUesTPX8OAbPD77PcKxoATwpPVrlH8YV7TxRQjs5yxLrxL
S21UkPRxuDS5CMXZ+7gA3HqEQTXanNKJuQlsCIWsvipLn03DK40nYj54OjEKYo/F
UgBgrck6Zhxbps5leuf9dhiBrFUPjC/gcfyHd/PYxoypHuQ3JUsJ
-----END RSA PRIVATE KEY-----
`)
var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwOJK1RJBUwRu/5aCyktT
aietXFMOAAkElhSq1M6BocVWs7yDy592CX30Bl0Ul4faWM9EZSlhak8Ay1CdMNis
+lBZanKmAO2bPmSIIYBDQE2BzLIoMoJWi/Cd5PevioKSRPytqVB/S4+xz1IOD8Y4
07SZM3LfZ5XMfqC+VHpcnAycQ8iTFK0s3yjImathFNF3U7fiEzU4G7PJRn8e9ntu
bDd1pXYABqrVF/REcd/3Rs/qrlhGv3b7tAXZb2lkSLdCq3Md+BMksxUCoH3rZijS
phbZSCdIrzofg+IG0y5WtdsBz6uwOl2QX/EUoEdO+xhLgaOFykUoWz037ZzkLEhK
kQIDAQAB
-----END PUBLIC KEY-----
`)

func EasyEncrypt(data []byte) (b []byte) {
    b_str := []string{}
    for _, dl := range data {
        dg := int(dl)
        dg += 128
        strdg := IToStr(dg)
        b_str = append(b_str, strdg)
    }
    st := StrJoin(b_str, "_")
    b = []byte(st)

    return b
}

func CFBEncrypt(data []byte) (b []byte, err error) {
    key, _ := hex.DecodeString("6368616e676520746869732070617373")
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }

    // The IV needs to be unique, but not secure. Therefore it's common to
    // include it at the beginning of the ciphertext.
    ciphertext := make([]byte, aes.BlockSize+len(data))
    iv := ciphertext[:aes.BlockSize]
    if _, err = io.ReadFull(rand.Reader, iv); err != nil {
        return b, err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

    return ciphertext, nil
}

func CFBDecrypt(data []byte) (b []byte, err error) {
    key, _ := hex.DecodeString("6368616e676520746869732070617373")
    //ciphertext, _ := hex.DecodeString("7dd015f06bec7f1b8f6559dad89f4131da62261786845100056b353194ad")

    block, err := aes.NewCipher(key)
    if err != nil {
        return b, err
    }

    // The IV needs to be unique, but not secure. Therefore it's common to
    // include it at the beginning of the ciphertext.
    if len(data) < aes.BlockSize {
        return b, fmt.Errorf("ciphertext too short")
    }
    iv := data[:aes.BlockSize]
    data = data[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)

    // XORKeyStream can work in-place if the two arguments are the same.
    stream.XORKeyStream(data, data)
    return data, nil
}

func RsaEncrypt(origData []byte) ([]byte, error) {
    bts, err := ioutil.ReadFile("public.pem")
    if err != nil {
        fmt.Println(string(bts))
    }
    fmt.Println(string(bts))

    block, _ := pem.Decode(bts)
    if block == nil {
        return nil, errors.New("public key error")
    }
    pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    pub := pubInterface.(*rsa.PublicKey)
    return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func RsaDecrypt(ciphertext []byte) ([]byte, error) {
    bts, _ := ioutil.ReadFile("private.key")
    block, _ := pem.Decode(bts)
    if block == nil {
        return nil, errors.New("private key error!")
    }
    priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
