FROM python:3.13-alpine

WORKDIR /usr/src/app

COPY . .
RUN pip3 install --no-cache-dir -r requirements.txt

CMD [ "python", "./main.py" ]
