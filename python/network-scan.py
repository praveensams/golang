import nmap

nm=nmap.PortScanner()

print(nm.scan('www.google.com', '22-443'))
