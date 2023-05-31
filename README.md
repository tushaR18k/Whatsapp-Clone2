## WhatsAPP Clone
### This project is build using React and Golang.
* For storage Postgre SQL is used using GORM package
* Websocket is used for real time chatting
* Future additions I am working on - 
    * Adding the feature to send files
    * Recording audio messages and sending
    * Calling via audio/video calls

## Prerequisites
* Make sure `Postgres is running` and there exists a database by the name `whatsapp`.

## To run
* Clone the repo
    * `cd` to whatsapp-backend and run command `go run .`
    * In another terminal `cd` to whatsapp-frontend and run command `npm start`
    * A browser window will pop-up. Navigate to `http://localhost:3000/signup`