package main

import (
	"QuakeAPI/core"
	"QuakeAPI/log"
	"QuakeAPI/utils"
	"bytes"
	"fmt"
	"strings"
)

func main() {
	utils.PrintLogo()
	input := utils.ParseInput()
	quakeCore := core.Core{}

	if input.UserInfo {
		quakeCore.GetUserInfo(input.Key)
	}

	if len(input.Search) != 0 && strings.TrimSpace(input.Search) != "" {
		var results string
		buffer := bytes.Buffer{}
		buffer.WriteString("IP\tPort\tService\tDomain/Host\tTitle\tLocation\n")

		remaining := input.Total
		start := 0
		collected := 0

		for remaining > 0 {
			fetchSize := 100
			if remaining < 100 {
				fetchSize = remaining
			}

			count, result := quakeCore.GetServiceInfo(input.Key, input.Search, fetchSize, start)

			if result == "" || count == 0 {
				break
			}

			buffer.WriteString(result)
			start += count
			remaining -= count
			collected += count

			fmt.Printf("[+] Progress: Collected %d records...\n", collected)

			if count < fetchSize {
				break
			}
		}
		results = buffer.String()

		if collected == 0 {
			log.Log("No data written because query returned no results or request failed", log.ERROR)
			return
		}

		utils.WriteOutput(results, input.Output)
	}
}
