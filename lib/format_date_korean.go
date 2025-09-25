package lib

import "time"

func FormatDateKorean(dateString string) (string, error) {
	inputLayout := "2006-01-02"

	t, err := time.Parse(inputLayout, dateString)
	if err != nil {
		return "", err
	}

	outputLayout := "2006년 01월 02일"

	return t.Format(outputLayout), nil
}
