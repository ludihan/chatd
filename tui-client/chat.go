package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	viewport        viewport.Model
	messages        []string
	textarea        textarea.Model
	senderStyle     lipgloss.Style
	err             error
	conn            *amqp.Connection
	ch              *amqp.Channel
	exchange        string
	nome            string
	urlAmqp         string
	urlApi          string
	msgs            <-chan amqp.Delivery
	previousMessage string
}

func initialModel() model {
	if len(os.Args) < 5 {
		fmt.Println("Informe o seu nome, a exchange, o urlAmqp do amqp e o url da api")
		os.Exit(1)
	}

	nome := os.Args[1]
	exchange := os.Args[2]
	urlAmqp := os.Args[3]
	urlApi := os.Args[4]

	conn, err := amqp.Dial(urlAmqp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ch.ExchangeDeclare(
		exchange, // name
		"fanout", // type
		false,    // durable
		true,     // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ch.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ta := textarea.New()
	ta.Placeholder = "Mande uma mensagem"
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 30)
	vp.SetContent(fmt.Sprintf("Bem vindo à exchange \"%v\"", exchange))

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		conn:        conn,
		ch:          ch,
		exchange:    exchange,
		nome:        nome,
		urlAmqp:     urlAmqp,
		urlApi:      urlApi,
		msgs:        msgs,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	select {
	case newMessage := <-m.msgs:
		message := struct {
			Exchange string `json:"exchange"`
			Body     string `json:"body"`
			UserId   string `json:"userId"`
		}{}
		json.Unmarshal(newMessage.Body, &message)
		if message.Body == m.previousMessage && message.UserId == m.nome {
			// nada
		} else {
			m.messages = append(m.messages, m.senderStyle.Render(fmt.Sprintf("%v: ", message.UserId))+message.Body)
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	default:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				fmt.Println(m.textarea.Value())
				m.ch.Close()
				m.conn.Close()
				return m, tea.Quit
			case tea.KeyEnter:
				m.messages = append(m.messages, m.senderStyle.Render(fmt.Sprintf("%v: ", m.nome))+m.textarea.Value())
				m.previousMessage = m.textarea.Value()

				body, _ := json.Marshal(map[string]string{
					"exchange": m.exchange,
					"body":     m.textarea.Value(),
					"userId":   m.nome,
				})

				payload := bytes.NewBuffer(body)

				http.Post(m.urlApi, "application/json", payload)

				m.viewport.SetContent(strings.Join(m.messages, "\n"))
				m.textarea.Reset()
				m.viewport.GotoBottom()
			}

			// We handle errors just like any other message
		case errMsg:
			m.err = msg
			return m, nil
		}
	}
	return m, tea.Batch(tiCmd, vpCmd)

}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"
}
