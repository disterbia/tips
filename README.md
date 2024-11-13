
1. rds 생성 후 ec2에서 rds 엔드포인트로 접속해서 db 생성해줘야함.
2. rds 인바운드규칙 설정해도 vpc를 열어줘야 로컬 접속가능. 참조 : https://gksdudrb922.tistory.com/240

https 설정

sudo certbot certonly --standalone -d haruharulab.com 

sudo certbot certonly --standalone -d wellkinson.haruharulab.com 

sudo nano /etc/nginx/sites-available/default

server {
    listen 80;
    server_name haruharulab.com;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    location / {
        return 301 https://$host$request_uri;
    }
}

server {
    listen 443 ssl;
    server_name haruharulab.com;

    ssl_certificate /etc/letsencrypt/live/haruharulab.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/haruharulab.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;
    ssl_ciphers HIGH:!aNULL:!MD5;

    client_max_body_size 20M;

    location / {
        proxy_pass http://localhost:40000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
    listen 80;
    server_name wellkinson.haruharulab.com;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    location / {
        return 301 https://$host$request_uri;
    }
}

server {
    listen 443 ssl;
    server_name wellkinson.haruharulab.com;

    ssl_certificate /etc/letsencrypt/live/wellkinson.haruharulab.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/wellkinson.haruharulab.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;
    ssl_ciphers HIGH:!aNULL:!MD5;

    client_max_body_size 20M;   

    location / {
        proxy_pass http://localhost:50000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}


인증서 갱신 : 

sudo systemctl stop nginx
sudo certbot renew --dry-run
sudo nginx -t
sudo systemctl restart nginx
sudo certbot renew --webroot -w /var/www/html

랜딩페이지 :
 sudo chown -R ubuntu:ubuntu /var/www/flutterweb
 $ sudo chmod -R 755 /var/www/flutterweb

 scp -i /Users/admin/Desktop/wellkinson.pem -r web/* ubuntu@43.203.141.35:/var/www/flutterweb

server {
    listen 80;
    server_name wellkinson.kr;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    location / {
        return 301 https://$host$request_uri;
    }
}

server {
    listen 443 ssl;
    server_name wellkinson.kr;

    ssl_certificate /etc/letsencrypt/live/wellkinson.kr/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/wellkinson.kr/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;
    ssl_ciphers HIGH:!aNULL:!MD5;

    root /var/www/flutterweb;
    index index.html;

    location / {
        try_files $uri /index.html;
    }
}


