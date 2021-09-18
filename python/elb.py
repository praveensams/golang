import boto3
elb = boto3.client('elb')
dns_records = []
lbs=elb.describe_load_balancers()
for lb in lbs["LoadBalancerDescriptions"]:
    print(lb["DNSName"])