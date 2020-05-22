package main

import (
        "fmt"
        "strings"
        "time"
        "os/exec"
        "encoding/json"
        "go.uber.org/zap"

        "github.com/yanzay/tbot/v2"
)

var (
	// for telegram(Change Required)
	Bot       *tbot.Server = tbot.New("1234~:AAAAAAAAAAA~")
        ChatID string = "1234~"

	Paths string = "hub_chain"
	SrcChainID string = "hub"
	DstChainID string = "chain"
	ClientID string = "srchclientaf"
)

type Result struct {
        Height string
        Txhash string
        RawLog string `json:"raw_log"`
}

func main() {

        var resRly Result
        var resRlyRe Result
        log,_ := zap.NewDevelopment()

        go Start()

        for {
                Send("```\n" +"Start] hub_" +DstChainID +"(4m)\n\n```")

                outRly, err := Command("rly tx raw update-client " +SrcChainID +" " +DstChainID +" " +ClientID)
                json.Unmarshal(outRly, &resRly)

                fmt.Println("err: ", err)


                if err != nil {
                        Send("!! Update command false")
                        log.Fatal("Update_rly Client", zap.Bool("Success", false), zap.String("err", string(outRly),))
                } else {
                        Send("```" + "\n> result: " +string(outRly) +"\n\n> Tx Hash: "+resRly.Txhash +"```")
                        log.Info("Update_rly Client", zap.Bool("Success", true), zap.String("err", "nil"), zap.String("tx", resRly.Txhash),)

                        k:=0
                        for {
                                if strings.Contains(string(outRly), "codespace") {
                                        Send("``` \n!!! Codespace !!!!\n```")
                                        outRlyRe, _ := Command("rly tx raw update-client " +SrcChainID +" " +DstChainID +" " +ClientID)
                                        json.Unmarshal(outRlyRe, &resRlyRe)
                                        Send("> Retry_Chain ID: hub_ " +DstChainID +" -\n\n" +"```" + "\n> result: " +string(outRlyRe) +"\n\n> Tx Hash: "+resRlyRe.Txhash +"```")

                                        if k ==2 {
                                                Send("``` \n!!!!!!!!! Requires restart !!!!!!!!!!\n```")
                                                break;
                                        }
                                } else {
                                        break;
                                }

                                k++
                                time.Sleep(time.Second *5)
                        }
                }

                for {
                        outTmp, _ := Command("rly tx rly "+Paths)
                        fmt.Printf("outTmp: %s\n", strings.TrimSpace(string(outTmp)))

                        if strings.Contains(string(outTmp), "No packets to relay") {
                                Send("```\n" +fmt.Sprint("> Relay Success! \n" +fmt.Sprint("Relay result: ", string(outTmp))) +"\n```")
                                break;
                        } else {
                                outTmp2, _ := Command("rly tx relay " +Paths)
                                Send("```\n" +fmt.Sprint("> Relay False! \n" +fmt.Sprint("Relay retry result: ", string(outTmp2)) +"\n```"))
                        }

                        time.Sleep(time.Second * 3)
                }

                // balance check
                outBalance, _ := Command("rly q balance " +SrcChainID)
                Send("```\n" +"End] Balance: " +string(outBalance) +"\n\n```")
                Send("-------------------------------------------------")

                // 4m
                time.Sleep(time.Second * 480)

        }
}

func Run() {

        fmt.Println("telegram Run")
//      responses := map[string]func() string{
        _ = map[string]func() string{
                "/start": func() string { return "Nice to meet you!" },
                "/hi":    func() string { return "Hi!" },
                "/now":   func() string { return fmt.Sprintf("%v", time.Now()) },
                "/msg":   func() string { return fmt.Sprintf("%v", "what?") },
        }
/*
        Bot.HandleMessage(".", func(m *tbot.Message) {
                resp, ok := responses[m.Text]
                if !ok {
                        Bot.Client().SendMessage(m.Chat.ID, "I didn't understand you.")
                        return
                }
                Bot.Client().SendMessage(m.Chat.ID, resp())
        })
*/
//      log.Fatal(Bot.Start())
        Bot.Start()

}


func Start() {
        Bot.Start()
}


func Send(str string) {

        Bot.Client().SendMessage(ChatID, str, tbot.OptParseModeMarkdown)

}



func Command(str string) ([]byte, error) {
        cmd := str
        out, err := exec.Command("/bin/bash", "-c", cmd).Output()
//      fmt.Println(cmd)

        return out, err
}

