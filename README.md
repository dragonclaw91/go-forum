# go forum
might need to revert the changes for the ionic serve in a bit 

Update the CLI: npm install -g @ionic/cli@latest

Revert the Config: Change "type": "custom" back to "type": "angular" in ionic.config.json.

 If it works without the error, you delete the ionic:serve line from your package.json.