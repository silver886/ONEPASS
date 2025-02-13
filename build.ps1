$projectName = 'ONEPASS'
$projectVersion = '1.0.0'

$gui = $false
$pack = $false

Push-Location $PSScriptRoot

$timestamp = [Xml.XmlConvert]::ToString((Get-Date)) -replace '(\+|-|:|\.)', '' -replace '(Z|\d\d\d\d$)', ''

$outputName = "$($projectName)_$($projectVersion)_$timestamp"

goversioninfo -64 -o "$projectName.syso" "$projectName.json"

dotenvx get --overload --strict -f .\.env.mapleridge --format=eval | ForEach-Object {
    Invoke-Expression -Command "`$$_"
}

if ($gui) {
    if ($pack) {
        packr build -o "$outputName.exe" -ldflags="-H=windowsgui -s -w -X main.defaultCloudflaredVersion=$CLOUDFLARED_VERSION -X main.defaultCloudflareZeroOrganization=$CLOUDFLARE_ZERO_ORGANIZATION -X main.defaultRemote=$REMOTE -X main.defaultLocal=$LOCAL"
    }
    else {
        go build -o "$outputName.exe" -ldflags="-H=windowsgui -s -w -X main.defaultCloudflaredVersion=$CLOUDFLARED_VERSION -X main.defaultCloudflareZeroOrganization=$CLOUDFLARE_ZERO_ORGANIZATION -X main.defaultRemote=$REMOTE -X main.defaultLocal=$LOCAL"
    }
}
else {
    if ($pack) {
        packr build -o "$outputName.exe" -ldflags="-s -w -X main.defaultCloudflaredVersion=$CLOUDFLARED_VERSION -X main.defaultCloudflareZeroOrganization=$CLOUDFLARE_ZERO_ORGANIZATION -X main.defaultRemote=$REMOTE -X main.defaultLocal=$LOCAL"
    }
    else {
        go build -o "$outputName.exe" -ldflags="-s -w -X main.defaultCloudflaredVersion=$CLOUDFLARED_VERSION -X main.defaultCloudflareZeroOrganization=$CLOUDFLARE_ZERO_ORGANIZATION -X main.defaultRemote=$REMOTE -X main.defaultLocal=$LOCAL"
    }
}

Remove-Item -Path "$projectName.syso"

upx --ultra-brute "$outputName.exe"

Copy-Item "$outputName.exe" "$projectName.exe"
