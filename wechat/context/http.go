package context

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"golang.org/x/crypto/pkcs12"
)

// HTTPGet get 请求
func (ctx *Context) HTTPGet(uri string) ([]byte, error) {
	response, err := ctx.RestyClient.R().Get(uri)
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode())
	}
	return response.Body(), nil
}

// HTTPPost post 请求
func (ctx *Context) HTTPPost(uri string, data string) ([]byte, error) {
	response, err := ctx.RestyClient.R().SetBody([]byte(data)).Post(uri)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode())
	}
	return response.Body(), nil
}

// PostJSON post json 数据请求
func (ctx *Context) PostJSON(uri string, obj interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	jsonData = bytes.Replace(jsonData, []byte("\\u003c"), []byte("<"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u003e"), []byte(">"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)
	response, err := ctx.RestyClient.R().SetHeader("Content-responseType", "application/json;charset=utf-8").SetBody(jsonData).Post(uri)
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode())
	}
	return response.Body(), nil
}

// PostJSONWithRespContentType post json数据请求，且返回数据类型
func (ctx *Context) PostJSONWithRespContentType(uri string, obj interface{}) ([]byte, string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, "", err
	}

	jsonData = bytes.Replace(jsonData, []byte("\\u003c"), []byte("<"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u003e"), []byte(">"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)

	response, err := ctx.RestyClient.R().SetHeader("Content-responseType", "application/json;charset=utf-8").SetBody(jsonData).Post(uri)
	if err != nil {
		return nil, "", err
	}
	if response.StatusCode() != http.StatusOK {
		return nil, "", fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode())
	}
	contentType := response.Header().Get("Content-responseType")
	return response.Body(), contentType, err
}

// PostFile 上传文件
func (ctx *Context) PostFile(fieldname, filename, uri string) ([]byte, error) {
	fields := []MultipartFormField{
		{
			IsFile:    true,
			Fieldname: fieldname,
			Filename:  filename,
		},
	}
	return ctx.PostMultipartForm(fields, uri)
}

// MultipartFormField 保存文件或其他字段信息
type MultipartFormField struct {
	IsFile    bool
	Fieldname string
	Value     []byte
	Filename  string
}

// PostMultipartForm 上传文件或其他多个字段
func (ctx *Context) PostMultipartForm(fields []MultipartFormField, uri string) (respBody []byte, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for _, field := range fields {
		if field.IsFile {
			fileWriter, e := bodyWriter.CreateFormFile(field.Fieldname, field.Filename)
			if e != nil {
				err = fmt.Errorf("error writing to buffer , err=%v", e)
				return
			}

			fh, e := os.Open(field.Filename)
			if e != nil {
				err = fmt.Errorf("error opening file , err=%v", e)
				return
			}
			defer fh.Close()

			if _, err = io.Copy(fileWriter, fh); err != nil {
				return
			}
		} else {
			partWriter, e := bodyWriter.CreateFormField(field.Fieldname)
			if e != nil {
				err = e
				return
			}
			valueReader := bytes.NewReader(field.Value)
			if _, err = io.Copy(partWriter, valueReader); err != nil {
				return
			}
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	bodyByte, err := io.ReadAll(bodyBuf)
	if err != nil {
		return nil, err
	}
	response, err := ctx.RestyClient.R().SetHeader("Content-responseType", contentType).SetBody(bodyByte).Post(uri)
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		return nil, err
	}
	respBody = response.Body()
	return
}

// PostXML perform a HTTP/POST request with XML body
func (ctx *Context) PostXML(uri string, obj interface{}) ([]byte, error) {
	xmlData, err := xml.Marshal(obj)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(xmlData)
	response, err := http.Post(uri, "application/xml;charset=utf-8", body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return io.ReadAll(response.Body)
}

// httpWithTLS CA证书
func httpWithTLS(rootCa, key string) (*http.Client, error) {
	var client *http.Client
	certData, err := os.ReadFile(rootCa)
	if err != nil {
		return nil, fmt.Errorf("unable to find cert path=%s, error=%v", rootCa, err)
	}
	cert := pkcs12ToPem(certData, key)
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tr := &http.Transport{
		TLSClientConfig:    config,
		DisableCompression: true,
	}
	client = &http.Client{Transport: tr}
	return client, nil
}

// pkcs12ToPem 将Pkcs12转成Pem
func pkcs12ToPem(p12 []byte, password string) tls.Certificate {
	blocks, err := pkcs12.ToPEM(p12, password)
	defer func() {
		if x := recover(); x != nil {
			log.Print(x)
		}
	}()
	if err != nil {
		panic(err)
	}
	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}
	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		panic(err)
	}
	return cert
}

// PostXMLWithTLS perform a HTTP/POST request with XML body and TLS
func PostXMLWithTLS(uri string, obj interface{}, ca, key string) ([]byte, error) {
	xmlData, err := xml.Marshal(obj)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(xmlData)
	client, err := httpWithTLS(ca, key)
	if err != nil {
		return nil, err
	}
	response, err := client.Post(uri, "application/xml;charset=utf-8", body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return io.ReadAll(response.Body)
}
