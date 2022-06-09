package app

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// @@@SNIPSTART money-transfer-project-template-go-workflow
func TransferMoney(ctx workflow.Context, transferDetails TransferDetails) error {

	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    500,
	}
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failures by default, this is just an example.
		RetryPolicy: retrypolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	err := workflow.ExecuteActivity(ctx, Withdraw, transferDetails).Get(ctx, nil)
	if err != nil {
		return err
	}
	StatusMessage("Withdraw")

	err = workflow.ExecuteActivity(ctx, Deposit, transferDetails).Get(ctx, nil)
	if err != nil {
		return err
	}
	StatusMessage("Deposit")

	return nil
}

func StatusMessage(data string) {

	var file string = "successLog.txt"
	exist := fileExists(file)

	if !exist {
		os.Create(file)
	}

	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	if err != nil {
		log.Fatalln("error starting TransferMoney workflow", err)
	} else {

		msg := "Successfully " + data + " of money"
		data := fmt.Sprintf("%v: %v", time.Now().Format("2006-01-02 3:4:5 pm"), msg)

		_, err = f.WriteString(data)
		f.WriteString("\n")

		if err != nil {
			fmt.Println("File reading error", err)
			return
		}
		f.Close()

	}

}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// @@@SNIPEND
