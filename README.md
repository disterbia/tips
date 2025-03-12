
1. rds 생성 후 ec2에서 rds 엔드포인트로 접속해서 db 생성해줘야함.
2. rds 인바운드규칙 설정해도 vpc를 열어줘야 로컬 접속가능. 참조 : https://gksdudrb922.tistory.com/240

protoc --go_out=. --go-grpc_out=. email.proto

https 설정

sudo certbot certonly --standalone -d adapfit-plus.com

sudo certbot certonly --standalone -d wellkinson.haruharulab.com 

sudo nano /etc/nginx/sites-available/default

docker system prune -a -f

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
    server_name wellkinson.haruharulab.com; //지금안씀 구 웰킨슨

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


인증서 자동갱신 : 

1.기존 certbot.timer 중지 (기본 설정 비활성화)
systemctl list-timers | grep certbot
sudo systemctl stop certbot.timer
sudo systemctl disable certbot.timer
2.certbot 갱신용 스크립트 만들기
sudo nano /usr/local/bin/certbot-renew.sh
#!/bin/bash

# 1. Nginx 중지
echo "Stopping Nginx..."
sudo systemctl stop nginx

# 2. Certbot으로 인증서 갱신
echo "Renewing SSL certificates..."
sudo certbot renew

# 3. Nginx 다시 시작
echo "Starting Nginx..."
sudo systemctl start nginx

echo "SSL renewal process completed."

저장후 권한추가
sudo chmod +x /usr/local/bin/certbot-renew.sh

3. 크론탭에 자동 실행 추가 (매일 새벽 1시UTC 실행)
sudo crontab -e
0 1 * * * /usr/local/bin/certbot-renew.sh >> /var/log/certbot-renew.log 2>&1

4. 로그확인:cat /var/log/certbot-renew.log

랜딩페이지 :

 sudo chmod -R 755 /var/www/wellkinson
sudo chmod 644 /var/www/kldga/sitemap.xml
 sudo chown -R www-data:www-data /var/www/wellkinson
 sudo chown -R ubuntu:ubuntu /var/www/wellkinson
 sudo chown -R www-data:www-data /var/www/wellkinson
 sudo chmod 644 /var/www/kldga/robots.txt
 sudo chown www-data:www-data /var/www/kldga/sitemap.xml


sudo rm -rf /var/www/wellkinson
 scp -i /Users/admin/Desktop/wellkinson.pem -r wellkinson_web/* ubuntu@54.180.236.32:/var/www/wellkinson

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

    root /var/www/wellkinson;
    index index.html;

    location / {
        try_files $uri /index.html;
    }
     # Sitemap 설정
    location /sitemap.xml {
        root /var/www/wellkinson;
    }

    # Robots.txt 설정
    location /robots.txt {
        root /var/www/wellkinson;
    }

    # Gzip 압축 설정
    gzip on;
    gzip_types text/plain text/css application/javascript application/json image/svg+xml;
    gzip_min_length 256;

    # 정적 파일 캐싱 설정
    location ~* \.(?:ico|css|js|gif|jpe?g|png|woff2?|eot|ttf|svg|map)$ {
        expires 6M;
        access_log off;
    }
}

server {
    listen 80;
    server_name kawa-official.org;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    location / {
        return 301 https://$host$request_uri;
    }
}

server {
    listen 443 ssl;
    server_name kawa-official.org;

    ssl_certificate /etc/letsencrypt/live/kawa-official.org/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/kawa-official.org/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;
    ssl_ciphers HIGH:!aNULL:!MD5;

    root /var/www/kawa;
    index index.html;

    location / {
        try_files $uri /index.html;
    }
     # Sitemap 설정
    location /sitemap.xml {
        root /var/www/kawa;
    }

    # Robots.txt 설정
    location /robots.txt {
        root /var/www/kawa;
    }

    # Gzip 압축 설정
    gzip on;
    gzip_types text/plain text/css application/javascript application/json image/svg+xml;
    gzip_min_length 256;

    # 정적 파일 캐싱 설정
    location ~* \.(?:ico|css|js|gif|jpe?g|png|woff2?|eot|ttf|svg|map)$ {
        expires 6M;
        access_log off;
    }
}

server {
    listen 80;
    server_name kldga.kr;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    location / {
        return 301 https://$host$request_uri;
    }
    
}

server {
    listen 443 ssl;
    server_name kldga.kr;

    ssl_certificate /etc/letsencrypt/live/kldga.kr/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/kldga.kr/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;
    ssl_ciphers HIGH:!aNULL:!MD5;

    root /var/www/kldga;
    index index.html;

    # Flutter Web 기본 라우팅 설정
    location / {
        try_files $uri /index.html;
    }

    # Sitemap 설정
    location /sitemap.xml {
        root /var/www/kldga;
    }

    # Robots.txt 설정
    location /robots.txt {
        root /var/www/kldga;
    }

    # Gzip 압축 설정
    gzip on;
    gzip_types text/plain text/css application/javascript application/json image/svg+xml;
    gzip_min_length 256;

    # 정적 파일 캐싱 설정
    location ~* \.(?:ico|css|js|gif|jpe?g|png|woff2?|eot|ttf|svg|map)$ {
        expires 6M;
        access_log off;
    }
}

server {
    listen 80;
    server_name adapfit-plus.com;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    location / {
        return 301 https://$host$request_uri;
    }
    
}

server {
    listen 443 ssl;
    server_name adapfit-plus.com;

    ssl_certificate /etc/letsencrypt/live/adapfit-plus.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/adapfit-plus.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;
    ssl_ciphers HIGH:!aNULL:!MD5;

    root /var/www/adapfit-plus;
    index index.html;

    # Flutter Web 기본 라우팅 설정
    location / {
        try_files $uri /index.html;
    }

    # Sitemap 설정
    location /sitemap.xml {
        root /var/www/adapfit-plus;
    }

    # Robots.txt 설정
    location /robots.txt {
        root /var/www/adapfit-plus;
    }

    # Gzip 압축 설정
    gzip on;
    gzip_types text/plain text/css application/javascript application/json image/svg+xml;
    gzip_min_length 256;

    # 정적 파일 캐싱 설정
    location ~* \.(?:ico|css|js|gif|jpe?g|png|woff2?|eot|ttf|svg|map)$ {
        expires 6M;
        access_log off;
    }
}

