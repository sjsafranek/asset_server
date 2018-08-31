import requests

r = requests.post("http://localhost:1111/api/v1/asset", files={
    'uploadfile': open('test.jpg','rb')
})

print(r.text)
if 200 != r.status_code:
    exit()

asset_id = r.json()['data']['asset_id']

r = requests.get("http://localhost:1111/api/v1/asset/{0}".format(asset_id))
print(r.text)
if 200 != r.status_code:
    exit()

r = requests.delete("http://localhost:1111/api/v1/asset/{0}".format(asset_id))
print(r.text)
if 200 != r.status_code:
    exit()

r = requests.get("http://localhost:1111/api/v1/asset/{0}".format(asset_id))
print(r.text)
if 200 != r.status_code:
    exit()
