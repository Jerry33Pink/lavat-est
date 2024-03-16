package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/near/near-api-go/accounts"
    "github.com/near/near-api-go/near"
    "github.com/near/near-api-go/utils"
)

const rpcURL = "https://near-testnet.lava.build/lava-referer-d3c83feb-cf99-44f8-983c-b024389e04df"

func main() {
    near, err := near.NewConnection(rpcURL)
    if err != nil {
        log.Fatalf("Failed to create NEAR connection: %v", err)
    }

    var wallets []accounts.Account

    for i := 0; i < 10; i++ {
        wallet, err := near.AddAccount()
        if err != nil {
            log.Fatalf("Failed to create wallet: %v", err)
        }
        wallets = append(wallets, wallet)
    }

    for {
        if len(wallets) == 0 {
            fmt.Println("All wallets are empty. Exiting...")
            os.Exit(0)
        }

        for _, wallet := range wallets {
            amount := utils.NearToYocto("1")
            receiver := "your-receiver-account.near"
            err := sendTransaction(wallet, receiver, amount)
            if err != nil {
                log.Printf("Failed to send transaction from wallet %s: %v", wallet.AccountID, err)
            } else {
                fmt.Printf("Transaction sent from wallet %s to %s\n", wallet.AccountID, receiver)
            }
        }

        time.Sleep(2 * time.Minute)
    }
}

func sendTransaction(wallet accounts.Account, receiver string, amount string) error {
    transaction := near.NewTransaction(wallet.AccountID, receiver, amount)
    signature, err := wallet.Sign(transaction)
    if err != nil {
        return fmt.Errorf("failed to sign transaction: %v", err)
    }
    transaction.Signature = signature
    _, err = wallet.Provider.SendTransaction(context.Background(), transaction)
    if err != nil {
        return fmt.Errorf("failed to send transaction: %v", err)
    }
    return nil
}
