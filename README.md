
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

 sudo chmod -R 755 /var/www/wellkinson
sudo chmod 644 /var/www/wellkinson/sitemap.xml
 sudo chown -R www-data:www-data /var/www/wellkinson
 sudo chown -R ubuntu:ubuntu /var/www/wellkinson
 sudo chown -R www-data:www-data /var/www/wellkinson
 sudo chmod 644 /var/www/wellkinson/robots.txt
 sudo chown www-data:www-data /var/www/wellkinson/robots.txt


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

