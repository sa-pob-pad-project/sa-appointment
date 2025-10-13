FROM golang:1.24.4

WORKDIR /app
COPY . .

# สร้างไบนารีชื่อ app
RUN go build -o app .

# ปรับพอร์ตตามโค้ดคุณ
EXPOSE 8001

CMD ["./app"]
