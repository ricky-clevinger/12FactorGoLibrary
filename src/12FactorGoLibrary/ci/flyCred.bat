@echo off
echo.
echo.
echo ** This batch assumes the first parameter is the name of the target and the pipeline and the second parameter is the pipeline path, the third parameter is credentials.yml.
echo.
fly --target %1 login --concourse-url http://192.168.100.4:8080
fly set-pipeline --target %1 --config %2 --pipeline %1 --load-vars-from %3
fly --target %1 unpause-pipeline --pipeline %1