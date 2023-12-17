from django.shortcuts import render, redirect
from django.http import JsonResponse, HttpResponse
import socket
from datetime import datetime
import requests

def new_link_redirect(request):
    return redirect('new/')

def new_link(request):
    if request.method == 'POST':
        try:
            input_value = request.POST['inputValue']
        except KeyError:
            return JsonResponse({'error': 'Bad Request: Missing inputValue field'}, status=400)
        client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        client.connect(('127.0.0.1', 6379))
        request_data = f'post\n{input_value}\n'
        client.send(request_data.encode('utf-8'))
        response_data = client.recv(1024).decode('utf-8')
        client.close()
        return HttpResponse(response_data)
    else:
        return render(request, 'main/index.html')

def getting(URL, SourceIP, TimeInterval):
    data = {'URL': URL, 'SourceIP': SourceIP, 'TimeInterval': TimeInterval}
    requests.post('http://127.0.0.1:5000', json=data)
    print(data)

def get_link(request, value):
    if value == 'favicon.ico':
        return HttpResponse(status=204)
    elif value == 'create_link':
        return render(request, 'myapp/index.html')
    else:
        client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        client.connect(('127.0.0.1', 6379))
        time = datetime.now().time()
       
        request_data = f'get\nhttp://127.0.0.1:8080/{value}\n'
        
        client.send(request_data.encode('utf-8'))
        response_data = client.recv(1024).decode('utf-8')
        i = 0
        if i < 1:
            i += 1
            getting(response_data + "(" + str(value) + ")", request.get_host(), str(time)[:5])
        client.close()

        return redirect(response_data)