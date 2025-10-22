# terraform/opentofu module for Azure IP addresses

I want to have an opentofu module that provides all the IP addresses Azure uses.
The output should provide different categories based on the parameters the input json provides.
The data we base on is https://download.microsoft.com/download/7/1/d/71d86715-5596-4529-9b13-da13a5de5b63/ServiceTags_Public_20251013.json.

Categories we want to have:
- All IP addresses
- IP range by "id", the ID can be eg AzurePortal or AzurePortal.SwedenCentral

For all categories, we want to have the option to get only IPv4, only IPv6 or both.

## Updates

As this data changes frequently, we should have a way to update the data in the module.
It should run once per night and check if thers is a new version of the json file.
If there is a new version, it should update the data in the module.
It should run as GitHub Action.
Changes are committed to the main branch directly.
The page which points to the latest version is https://www.microsoft.com/en-us/download/details.aspx?id=56519.
Make the script check for a link in a href that starts with "https://download.microsoft.com/download/" and ends with ".json".
If there are non or multiple links like that, the script should fail.
This way the maintainer gets notified that something is wrong.
Use a script in bash and call in the GitHub Action.

Preferably, there can be found a permalink to the latest version of the json file, so we don't have to scrape the download page.

## Technical details

I want to have no input.
I think we can have some yaml files that are generated from the json file.
Those yaml files are per "id" and contain the IP ranges for that id.
Then terraform can read those yaml files and provide the output based on the input parameters.

Use go for the script that generates the yaml files from the json file.

Use the change number and a diff between the old and new json file to determine if there are changes.

Use go.yaml.in/yaml/v3 and go 1.25.
For logging use charmbrace/log.

Do the download, check and generation in one go program.
Place the go program in useful different files.

After an update, the module needs to publish a tag in the formar "vYYYY.MM.DD" where YYYY is the year, MM the month and DD the day of the update.
