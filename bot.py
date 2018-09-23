import requests

base_url = 'https://api.groupme.com/v3'
params = {
    'token': '<insert your token here, do not add it to git>',
}
data = {
    'bot_id': 'c1d5a8381af0c8b21fe3bcb2b6',
    'text': 'На здоровье',
}
resp = requests.post(base_url + '/bots/post', params=params, data=data)
print(resp)
