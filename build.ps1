$projectName = 'ONEPASS'
$projectVersion = '1.1.2'
$projectEnv = 'mapleridge'

$gui = $false
$pack = $false

Push-Location $PSScriptRoot

$timestamp = [Xml.XmlConvert]::ToString((Get-Date)) -replace '(\+|-|:|\.)', '' -replace '(Z|\d\d\d\d$)', ''
$outputName = "$($projectName)_$($projectVersion)_$($timestamp)_{{OS}}_{{ARCH}}"

dotenvx get --overload --strict "--env-file=.env.$projectEnv" --format=eval | ForEach-Object {
    Invoke-Expression -Command "`$$_"
}

$cc = 'go'
$cflags = ''
$ldflags = "-s -w -X main.toolName=$projectName -X main.toolVersion=$projectVersion -X main.defaultCloudflaredVersion=$CLOUDFLARED_VERSION -X main.defaultRemote=$REMOTE -X main.defaultLocal=$LOCAL"
if ($gui) {
    $ldflags = "$ldflags -H=windowsgui $ldflags"
}
if ($pack) {
    $cc = 'packr'
}

& {
    $os = 'windows'
    $arch = 'amd64'
    $tempOutputName = $outputName.Replace('{{OS}}', $os).Replace('{{ARCH}}', $arch)
    $cflags = "-o=$tempOutputName.exe $cflags"

    goversioninfo -64 -o "$projectName.syso" "$projectName.json"
    env "GOOS=$os" "GOARCH=$arch" "$cc" build "$cflags" "-ldflags=$ldflags"

    Remove-Item -Path "$projectName.syso"
    upx --ultra-brute "$tempOutputName.exe"
    Copy-Item "$tempOutputName.exe" "$($tempOutputName.Replace("_$timestamp", '')).exe"
}

& {
    $os = 'darwin'
    $arch = 'arm64'
    $tempOutputName = $outputName.Replace('{{OS}}', $os).Replace('{{ARCH}}', $arch)
    $cflags = "-o=$tempOutputName $cflags"

    env "GOOS=$os" "GOARCH=$arch" "$cc" build "$cflags" "-ldflags=$ldflags"

    Copy-Item "$tempOutputName" "$($tempOutputName.Replace("_$timestamp", ''))"
}
