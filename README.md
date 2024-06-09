# Balance Sweep

Modern bank accounts like Monzo offer the ability to 'round-up' transactions and move pennies to a savings pot/account. Currently this only works with card transactions, meaning bank transfers, P2P payments, direct debits and more do not work with round-ups.

To solve this, this program aims to sweep pennies from a Monzo current account balance into a designated savings pot at the end of each day.

## Prerequisites

Before running the program, you'll need to register a new API client on the [Monzo Developer Portal](https://developers.monzo.com). 

To register a new API client, log inot the Monzo Devloper Portal (remember to approve the login via your Monzo app) and click "New OAuth Client".

Then provide the following details for your OAuth client (Logo URL can remain blank):

```
Name: Balance Sweep
Redirect URL: http://127.0.0.1:21234/callback
Description: Balance Sweep Application
Confidentiality: True
```

Once the client is registered you will recieve a Client ID and a Client Secert. Make a note of these!

## How to use

Clone the repository (I like using the [GitHub CLI](https://cli.github.com/) for this)
```bash
gh repo clone skyth3r/balance-sweep
```

Install dependencies
```bash
go mod tidy
```

Set Client ID and Client Secert in environment variables
```bash
export MONZO_CLIENT_ID=YOUR_CLIENT_ID_HERE

export MONZO_CLIENT_SECRET=YOUR_CLIENT_SECRET_HERE
```

Run the program
```bash
go run ./
```

## Expected results

The first time this code is run, the client will start the OAuth flow, and attempt to open a browser with the login URL. On the login page, type in your email address linked to your personal Monzo account and then click the link sent to your email address and go back to the app.

You will then be prompted to open the Monzo app and grant access to the app by clicking "Allow access to your data". This process is related to Strong Customer Authentication. Once access has been granted via the Monzo app, go back to the app and press the [Enter] key to continue.

The app then looks up your uk_retail account ID (personal current account ID), and checks if the balance needs to be rounded down. If it does then it moves the pennies to the pot with the name 'Savings'.

For future uses of the app, the access and refresh token will be retrieved from the system's keystore (e.g. on MacOS it woud be retrieved from Keychain). You will be prompted to allow the app to access to the system's keystore. As this access token grants access to your live Monzo bank account I would advice granting access with the 'Allow' option rather than the 'Always Allow' option. 

## Supported accounts

This app is intended for use with the Monzo Personal Current Account (uk_retail) but can be modified to work with other account types. 
