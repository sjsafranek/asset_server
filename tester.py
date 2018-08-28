import requests

r = requests.post("http://localhost:1111/api/v1/upload", files={
    'uploadfile': open('test.jpg','rb')
})

print(r.text)
