# AWS SAM

This directory contains configuration to deploy golinks with the AWS Servless Application Model.

## Deploying

```sh
make deploy
```

This will create Lambda Functions to:
- serve the homepage,
- create new Links, and
- redirect names to existing Links

It will also provision an API Gateway to direct API endpoints to the relevant functions.

Finally, it creates a DynamoDB table as storage for the Links.

## Custom Domain Name

The following sections provides the steps for using a custom public domain within AWS.
This can be followed after you have successfully deployed golinks.

1. Route 53 
   1. Provision or use an existing AWS Route53 domain.
2. AWS Certificate Manager
   1. Request or use an existing public certificate for the domain you wish to use. 
   2. Ensure you click "Create records in Route 53".
3. API Gateway
   1. Use the sidebar and select "Custom domain names".
   2. "Create" a domain name with the FQDN. 
   3. Select "API Mappings" and "Configure API mappings". From the dropdown, select the golinks API that was deployed.
   Choose stage "Prod", do not specify a Path, and "Save".
4. Go back to Route 53
   1. Select your "Hosted Zone" and click "Create record".
   2. Choose the "Simple routing" policy.
   3. Define a simple record with the subdomain of your choice.
   4. Under "Value/Route traffic to", choose "Alias to API Gateway API", select the region, and search for the custom domain name created in the API Gateway.

Your deployment should be available at the domain momentarily.