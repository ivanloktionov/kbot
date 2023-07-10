# Overview 

Link monitor telebot - telegram bot that receives URL and starts monitoring it if the status code or content of the webpage has been changed - a user will be notified about the changes.

# Getting Started

## Download via GitHub

1. Upload the telebot using the repository https://github.com/ivanloktionov/kbot
2. Add telebot token using read -s TELE_TOKEN
3. Run the binary file: ./kbot start

### Note: using this method some dependencies should be installed

## Install via docker

1. Add telebot token using read -s TELE_TOKEN
2. Docker run: docker run -e "`echo -n TELE_TOKEN=$TELE_TOKEN`" ghcr.io/ivanloktionov/kbot:v4.0.0-f2c91c1-linux-amd64

## Example of user interaction with telebot

![IMG_5335](https://github.com/ivanloktionov/kbot/assets/71848058/10cab36b-5635-40aa-b1ca-a5d546e2e4fd)

## General guide for telegram users:
In order to start monitoring a link, you can send the following message:
/s URL

To stop receiving alerts, you can use the command:
/s stop
