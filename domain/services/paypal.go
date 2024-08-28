package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/invoice"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/purchaseorder"
	"github.com/mlbautomation/ProyectoEMLB/model"
)

const (
	ExpectedVerification = "SUCCESS"
	ExpectedStatus       = "completed"
)

const (
	EventTypeProduct = "PAYMENT.CAPTURE.COMPLETED"
)

type PayPal struct {
	servicePurchaseOrder purchaseorder.ServicePaypal
	serviceInvoice       invoice.ServicePaypal
}

func NewPaypal(pspp purchaseorder.ServicePaypal, ispp invoice.ServicePaypal) PayPal {
	return PayPal{
		servicePurchaseOrder: pspp,
		serviceInvoice:       ispp,
	}
}

func (pp PayPal) ProcessRequest(header http.Header, body []byte) error {
	//Recibimos la información en nuestra URL que paypal debería consumir para notificarme un webhook
	payPalRequestValidator, payPalRequestData, err := pp.parsePayPalRequestAndData(header, body)
	if err != nil {
		errMsg := fmt.Errorf("%s %w", "pp.parsePayPalRequest()", err)
		log.Println(errMsg)
		return errMsg
	}

	//Le pedimos a paypal que valide lo que hemos hocho en el paso anterior
	err = pp.validate(&payPalRequestValidator)
	if err != nil {
		log.Println(err)
		return err
	}

	//Si no devuelve error, es decir paypal a verificado que todo este bien, procesamos el pago
	err = pp.processPayment(&payPalRequestData)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (pp PayPal) parsePayPalRequestAndData(headers http.Header, body []byte) (model.PayPalRequestValidator, model.PayPalRequestData, error) {
	data := model.PayPalRequestData{}

	//Vuelco la información del body en la estructura model.PayPalRequestData{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return model.PayPalRequestValidator{}, model.PayPalRequestData{}, fmt.Errorf("%s %w", "json.Unmarshal()", err)
	}

	//Verifico si el EventType string `json:"event_type" es PAYMENT.CAPTURE.COMPLETED
	if data.EventType != EventTypeProduct {
		return model.PayPalRequestValidator{}, model.PayPalRequestData{}, fmt.Errorf("the event_type %q is not allowed", data.EventType)
	}

	return model.PayPalRequestValidator{
		//Lo obtengo del headers
		AuthAlgo:         headers.Get("Paypal-Auth-Algo"),
		CertURL:          headers.Get("Paypal-Cert-Url"),
		TransmissionID:   headers.Get("Paypal-Transmission-Id"),
		TransmissionSig:  headers.Get("Paypal-Transmission-Sig"),
		TransmissionTime: headers.Get("Paypal-Transmission-Time"),
		//El body es el webhook
		WebhookEvent: body,
		//El webhook ID en las variables de enviorement
		WebhookID: os.Getenv("WEBHOOK_ID"),
	}, data, nil
}

func (pp PayPal) validate(p *model.PayPalRequestValidator) error {

	//convertimos el modelo al formato JSON
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	//hacemos un Request, con el método post, al URL de validadción
	request, err := http.NewRequest(http.MethodPost, os.Getenv("VALIDATION_URL"), bytes.NewReader(data))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	//Crea un header con authorization tipo básica, y uni el username : password
	request.SetBasicAuth(os.Getenv("CLIENT_ID"), os.Getenv("SECRET_ID"))

	//Hacemos la petición a paypal y recibimos la respuesta
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(response)

	//Aquí trabajamos con la respuesta de paypal
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	//Si no es OK = 200
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("PayPal response with status code %d, body: %s", response.StatusCode, string(body))
	}

	//Si fue correcta, leemos el body
	bodyMap := make(map[string]string)
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return err
	}

	//pero nos va a interesar solo el estado de verificación SUCCESS
	if bodyMap["verification_status"] != ExpectedVerification {
		return fmt.Errorf("verification status is %s", bodyMap["verification_status"])
	}

	return nil
}

func (pp PayPal) processPayment(data *model.PayPalRequestData) error {

	//Verificamos el estado de ese webhook, en pago: completed
	if !strings.EqualFold(data.Resource.Status, ExpectedStatus) {
		return fmt.Errorf("el estado de la transacción: %q no es el estado esperado: %q", data.Resource.Status, ExpectedStatus)
	}

	//Leemo el custom ID que es la orden de compra para nosotros
	ID, err := uuid.Parse(data.Resource.CustomID)
	if err != nil {
		return fmt.Errorf("%s %w", "uuid.Parse()", err)
	}

	//Busca la orden de compra por ID
	order, err := pp.servicePurchaseOrder.GetByID(ID)
	if err != nil {
		return fmt.Errorf("%s %w", "pp.servicePurchaseOrder.GetByID(ID)", err)
	}

	//Convertimos el valor del webhook que esta en string a punto flotante
	value, err := strconv.ParseFloat(data.Resource.Amount.Value, 64)
	if err != nil {
		return fmt.Errorf("%s %w", "strconv.ParseFloat(data.Resource.Amount.Value, 64)", err)
	}
	value = math.Floor(value*100) / 100

	//Convertimos el valor del webhook que esta en string a punto flotante
	totalAmount := pp.servicePurchaseOrder.TotalAmount(order)
	totalAmount = math.Floor(totalAmount*100) / 100

	/*
		x := 12.3456
		fmt.Println(math.Floor(x*100) / 100) // 12.34 (round down)
		fmt.Println(math.Round(x*100) / 100) // 12.35 (round to nearest)
		fmt.Println(math.Ceil(x*100) / 100)  // 12.35 (round up)
	*/

	//Comparamos el webhook con la orden de compra
	if totalAmount != value {
		return fmt.Errorf("el valor recibido: %0.2f, es diferente al valor esperado %0.2f", value, totalAmount)
	}

	//Creamos la factura con referencia a la orden de compra
	return pp.serviceInvoice.Create(&order)
}
