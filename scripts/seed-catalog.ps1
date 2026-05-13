# PCDoma catalog seed script.
# Logs in as admin, wipes existing PCs/peripherals/setups, and re-creates
# a curated catalog of 13 PCs (with images), 12 peripherals and 4 setups.
#
# Usage:   powershell -ExecutionPolicy Bypass -File scripts\seed-catalog.ps1
# Assumes: docker compose is up; gateway on http://localhost:8090

$ErrorActionPreference = "Stop"
$Gateway              = "http://localhost:8090/api/v1"
$AdminEmail           = "test@pcdoma.kz"
$AdminPassword        = "password123"

# ── Helpers ─────────────────────────────────────────────────────────
$script:Token = $null
function Invoke-Api {
    param([string]$Method, [string]$Path, $Body)
    $h = @{ "Content-Type" = "application/json; charset=utf-8" }
    if ($script:Token) { $h["Authorization"] = "Bearer $($script:Token)" }
    $callArgs = @{ Method = $Method; Uri = "$Gateway$Path"; Headers = $h }
    if ($null -ne $Body) {
        $json = $Body | ConvertTo-Json -Depth 10 -Compress
        $callArgs["Body"] = [System.Text.Encoding]::UTF8.GetBytes($json)
    }
    return Invoke-RestMethod @callArgs
}

# ── 1. Login ────────────────────────────────────────────────────────
Write-Host "→ Login as $AdminEmail" -ForegroundColor Cyan
$auth = Invoke-Api -Method POST -Path "/auth/login" -Body @{
    email    = $AdminEmail
    password = $AdminPassword
}
$script:Token = $auth.tokens.access_token
if (-not $script:Token) { throw "Login succeeded but no access_token in response" }
Write-Host "  ✓ token acquired" -ForegroundColor Green

# ── 2. Wipe existing catalog ────────────────────────────────────────
Write-Host "→ Wiping existing catalog" -ForegroundColor Cyan
try { $existingSetups = (Invoke-Api GET "/catalog/setups" $null).data } catch { $existingSetups = @() }
foreach ($s in $existingSetups) {
    try { Invoke-Api DELETE "/catalog/setups/$($s.id)" $null | Out-Null } catch {}
}
try { $existingPeripherals = (Invoke-Api GET "/catalog/peripherals" $null).data } catch { $existingPeripherals = @() }
foreach ($p in $existingPeripherals) {
    try { Invoke-Api DELETE "/catalog/peripherals/$($p.id)" $null | Out-Null } catch {}
}
try { $existingPCs = (Invoke-Api GET "/catalog/pcs?limit=100" $null).data } catch { $existingPCs = @() }
foreach ($pc in $existingPCs) {
    try { Invoke-Api DELETE "/catalog/pcs/$($pc.id)" $null | Out-Null } catch {}
}
Write-Host "  ✓ wiped $($existingPCs.Count) PCs, $($existingPeripherals.Count) peripherals, $($existingSetups.Count) setups" -ForegroundColor Green

# ── 3. Locations (3 districts of Astana) ────────────────────────────
$LocEsil      = @{ district = "Есиль";    address = "ул. Кунаева 14";       lat = 51.1284; lng = 71.4310 }
$LocAlmaty    = @{ district = "Алматы";   address = "пр. Республики 52";    lat = 51.1605; lng = 71.4704 }
$LocSaryarka  = @{ district = "Сарыарка"; address = "ул. Иманова 19";       lat = 51.1773; lng = 71.4491 }

# ── 4. PCs (13 builds) ──────────────────────────────────────────────
$PCs = @(
    @{ name = "Titan RGB Pro";          img = "titan_rgb_black.webp";       loc = $LocEsil;
       cpu = "Intel Core i9-14900K";    gpu = "NVIDIA RTX 4080 SUPER 16GB"; ram = "64GB DDR5 6400 MHz"; storage = "2TB NVMe Samsung 990 Pro";
       hour = 2800; day = 16000 },

    @{ name = "ABIX Phantom RTX 5080";  img = "abix_rtx5080.webp";          loc = $LocEsil;
       cpu = "Intel Core i7-14700K";    gpu = "NVIDIA RTX 5080 16GB";       ram = "32GB DDR5 6000 MHz"; storage = "2TB NVMe WD Black SN850";
       hour = 2500; day = 14500 },

    @{ name = "Orange Inferno";         img = "orange_inferno.webp";        loc = $LocAlmaty;
       cpu = "AMD Ryzen 7 7800X3D";     gpu = "NVIDIA RTX 4070 Ti SUPER 16GB"; ram = "32GB DDR5 6000 MHz"; storage = "1TB NVMe Kingston KC3000";
       hour = 2200; day = 13000 },

    @{ name = "Compway Violet";         img = "compway_violet.webp";        loc = $LocAlmaty;
       cpu = "Intel Core i5-13600KF";   gpu = "NVIDIA RTX 4070 12GB";       ram = "32GB DDR5 5600 MHz"; storage = "1TB NVMe Samsung 980 Pro";
       hour = 1800; day = 10500 },

    @{ name = "ABIX Storm RTX 5060";    img = "abix_rtx5060_aio.webp";      loc = $LocSaryarka;
       cpu = "Intel Core i5-14400F";    gpu = "NVIDIA RTX 5060 8GB";        ram = "16GB DDR5 5600 MHz"; storage = "1TB NVMe Crucial P5 Plus";
       hour = 1500; day = 8800 },

    @{ name = "Snow Panorama";          img = "snow_panorama.webp";         loc = $LocEsil;
       cpu = "Intel Core i7-13700K";    gpu = "NVIDIA RTX 4070 SUPER 12GB"; ram = "32GB DDR5 6000 MHz"; storage = "1TB NVMe Samsung 990 Pro";
       hour = 2000; day = 11800 },

    @{ name = "CubeGaming Starter";     img = "cubegaming_starter.webp";    loc = $LocSaryarka;
       cpu = "AMD Ryzen 5 7600";        gpu = "NVIDIA RTX 4060 8GB";        ram = "16GB DDR5 5200 MHz"; storage = "512GB NVMe Kingston NV2";
       hour = 1100; day = 6500 },

    @{ name = "WhiteBox Office";        img = "whitebox_office.webp";       loc = $LocSaryarka;
       cpu = "Intel Core i5-13400";     gpu = "Intel UHD 730 (Integrated)"; ram = "16GB DDR4 3200 MHz"; storage = "512GB SSD Kingston A2000";
       hour = 600;  day = 3500 },

    @{ name = "Arctic Dual Chamber";    img = "arctic_dual_chamber.webp";   loc = $LocEsil;
       cpu = "Intel Core i9-13900K";    gpu = "NVIDIA RTX 4080 SUPER 16GB"; ram = "64GB DDR5 6400 MHz"; storage = "2TB NVMe Samsung 990 Pro";
       hour = 2700; day = 15500 },

    @{ name = "Glacier XL";             img = "glacier_xl.webp";            loc = $LocAlmaty;
       cpu = "AMD Ryzen 9 7950X3D";     gpu = "NVIDIA RTX 4090 24GB";       ram = "64GB DDR5 6400 MHz"; storage = "4TB NVMe Samsung 990 Pro";
       hour = 3500; day = 19500 },

    @{ name = "Aerocool Basic";         img = "aerocool_basic.webp";        loc = $LocSaryarka;
       cpu = "Intel Core i3-12100F";    gpu = "NVIDIA GTX 1660 SUPER 6GB";  ram = "16GB DDR4 3000 MHz"; storage = "256GB SSD Crucial BX500";
       hour = 500;  day = 3000 },

    @{ name = "NZXT H7 Frost";          img = "nzxt_h7_white.webp";         loc = $LocAlmaty;
       cpu = "AMD Ryzen 7 7700X";       gpu = "NVIDIA RTX 4070 Ti 12GB";    ram = "32GB DDR5 6000 MHz"; storage = "1TB NVMe Samsung 980 Pro";
       hour = 1900; day = 11200 },

    @{ name = "Thermaltake Core Liquid"; img = "thermaltake_core.webp";     loc = $LocEsil;
       cpu = "Intel Core i9-14900KS";   gpu = "NVIDIA RTX 4090 24GB";       ram = "128GB DDR5 6400 MHz"; storage = "4TB NVMe Samsung 990 Pro + 8TB HDD";
       hour = 4000; day = 22000 },

    # ── Приколы / Easter eggs ──────────────────────────────────────
    @{ name = "Картонный Геймерский Сетап"; img = "cardboard_gamer.jpg"; loc = $LocSaryarka;
       cpu = "Picasso Cardboard E-2024";    gpu = "Imagination Pro™ (нарисована маркером)";
       ram = "8 ГБ слов поддержки";          storage = "1 DVD-диск с CS 1.6";
       hour = 50; day = 200 },

    @{ name = "IBM Quantum System One 127q"; img = "ibm_quantum_127q.jpg"; loc = $LocEsil;
       cpu = "IBM Eagle r3 (127 кубитов)";   gpu = "Не нужна — мы считаем волновые функции";
       ram = "Superposition Cache";           storage = "Запутанные кубиты, ∞ ТБ";
       hour = 999999; day = 5000000 },

    @{ name = "IBM Quantum Heron 156q";       img = "ibm_quantum_heron.jpg"; loc = $LocEsil;
       cpu = "IBM Heron r2 (156 кубитов, fault-tolerant)";
       gpu = "Топологическая визуализация Bloch sphere";
       ram = "Coherence time 250 µs";         storage = "Холодильник на 15 миликельвин";
       hour = 1499999; day = 7500000 },

    @{ name = "Ретро-ПК MT-9017T";            img = "retro_mt9017.jpg"; loc = $LocSaryarka;
       cpu = "Intel Pentium III 800 МГц";    gpu = "S3 Trio 64V+ (4 МБ)";
       ram = "256 МБ SDRAM PC-133";           storage = "20 ГБ IDE Seagate ST320423A";
       hour = 100; day = 500 },

    @{ name = "Картошка-PC (Soviet edition)"; img = "potato_pc.jpg"; loc = $LocSaryarka;
       cpu = "Спаянный AMD K6-2 + картошка-радиатор";
       gpu = "Voodoo Banshee (если повезёт)";
       ram = "32 МБ EDO + 16 МБ из соседнего ПК";
       storage = "Дискета 3.5`" 1.44 МБ";
       hour = 30; day = 150 }
)

Write-Host "→ Creating $($PCs.Count) PCs" -ForegroundColor Cyan
$createdPCs = @{}
foreach ($p in $PCs) {
    $resp = Invoke-Api POST "/catalog/pcs" @{
        name           = $p.name
        specs          = @{ cpu = $p.cpu; gpu = $p.gpu; ram = $p.ram; storage = $p.storage }
        price_per_hour = $p.hour
        price_per_day  = $p.day
        location       = $p.loc
        images         = @("/images/pcs/$($p.img)")
    }
    $pcId = $resp.id
    $createdPCs[$p.name] = @{ id = $pcId; district = $p.loc.district }
    Write-Host "  + $($p.name) → $pcId" -ForegroundColor Gray
}

# ── 5. Peripherals (12) ─────────────────────────────────────────────
$Peripherals = @(
    # Есиль
    @{ name = "Razer DeathAdder V3";          type = "mouse";    loc = "Есиль";
       desc = "Лёгкая киберспортивная мышь, сенсор Focus Pro 30K, 64g.";
       specs = @{ sensor = "Focus Pro 30K"; dpi = "30000"; weight = "59g" };
       hour = 150; day = 700 },
    @{ name = "Keychron K2 Pro";              type = "keyboard"; loc = "Есиль";
       desc = "75% механическая клавиатура, hot-swap, Gateron G Pro Brown.";
       specs = @{ layout = "75%"; switch = "Gateron G Pro Brown"; backlight = "RGB" };
       hour = 250; day = 1200 },
    @{ name = "HyperX Cloud III";             type = "headset";  loc = "Есиль";
       desc = "Игровая гарнитура с микрофоном DTS Headphone:X.";
       specs = @{ driver = "53mm"; mic = "DTS"; weight = "320g" };
       hour = 200; day = 900 },
    @{ name = "LG UltraGear 27GP850 27`"";    type = "monitor";  loc = "Есиль";
       desc = "27`" IPS QHD 165 Гц, 1 мс, NVIDIA G-SYNC Compatible.";
       specs = @{ size = "27`""; resolution = "2560x1440"; refresh = "165Hz"; panel = "Nano IPS" };
       hour = 600; day = 3000 },

    # Алматы
    @{ name = "Logitech G Pro X Superlight 2"; type = "mouse";    loc = "Алматы";
       desc = "Беспроводная киберспортивная мышь, HERO 2 sensor, 60g.";
       specs = @{ sensor = "HERO 2"; dpi = "32000"; weight = "60g"; wireless = "LIGHTSPEED" };
       hour = 200; day = 950 },
    @{ name = "Razer Huntsman V3 Pro";        type = "keyboard"; loc = "Алматы";
       desc = "Аналоговая оптическая клавиатура, Rapid Trigger, 8000 Hz.";
       specs = @{ switch = "Razer Analog Optical"; polling = "8000Hz"; backlight = "Chroma RGB" };
       hour = 300; day = 1500 },
    @{ name = "SteelSeries Arctis Nova Pro";  type = "headset";  loc = "Алматы";
       desc = "Hi-Res наушники с активным шумоподавлением и DAC.";
       specs = @{ driver = "40mm"; mic = "ClearCast Gen 2"; anc = "yes"; dac = "GameDAC Gen 2" };
       hour = 350; day = 1700 },
    @{ name = "Samsung Odyssey G7 27`"";      type = "monitor";  loc = "Алматы";
       desc = "27`" 1000R изогнутый VA, QHD 240 Гц, 1 мс, HDR600.";
       specs = @{ size = "27`""; resolution = "2560x1440"; refresh = "240Hz"; panel = "VA Curved 1000R" };
       hour = 700; day = 3500 },

    # Сарыарка
    @{ name = "Logitech G502 HERO";           type = "mouse";    loc = "Сарыарка";
       desc = "Классическая игровая мышь с настраиваемыми грузами.";
       specs = @{ sensor = "HERO 25K"; dpi = "25600"; weight = "121g" };
       hour = 100; day = 500 },
    @{ name = "Logitech G413 SE";             type = "keyboard"; loc = "Сарыарка";
       desc = "Механическая клавиатура с тактильными свитчами.";
       specs = @{ switch = "Tactile"; backlight = "White LED"; layout = "Full" };
       hour = 150; day = 700 },
    @{ name = "Razer Kraken X";               type = "headset";  loc = "Сарыарка";
       desc = "Лёгкая гарнитура для длительных сессий, 250 г.";
       specs = @{ driver = "40mm"; weight = "250g"; surround = "7.1" };
       hour = 150; day = 700 },
    @{ name = "AOC 24G2 24`"";                type = "monitor";  loc = "Сарыарка";
       desc = "24`" IPS Full HD 144 Гц, 1 мс, FreeSync Premium.";
       specs = @{ size = "24`""; resolution = "1920x1080"; refresh = "144Hz"; panel = "IPS" };
       hour = 350; day = 1700 }
)

Write-Host "→ Creating $($Peripherals.Count) peripherals" -ForegroundColor Cyan
$createdPerifs = @{}
foreach ($p in $Peripherals) {
    $resp = Invoke-Api POST "/catalog/peripherals" @{
        name           = $p.name
        type           = $p.type
        description    = $p.desc
        specs          = $p.specs
        price_per_hour = $p.hour
        price_per_day  = $p.day
        location_id    = $p.loc
        images         = @()
    }
    $perifId = $resp.id
    $key = "$($p.loc)|$($p.type)|$($p.name)"
    $createdPerifs[$key] = $perifId
    Write-Host "  + $($p.name) [$($p.loc) / $($p.type)] → $perifId" -ForegroundColor Gray
}

function Get-Perif([string]$Loc, [string]$Type, [string]$Name) {
    return $createdPerifs["$Loc|$Type|$Name"]
}

# ── 6. Setups (4) ───────────────────────────────────────────────────
$Setups = @(
    @{ name = "Ultimate Gaming Setup"; category = "gaming"; discount = 10
       pc = "Glacier XL"
       desc = "Топовый гейминг для AAA-проектов и киберспорта на максималках."
       img  = "glacier_xl.webp"
       items = @(
         @{ loc = "Алматы"; type = "mouse";    name = "Logitech G Pro X Superlight 2" },
         @{ loc = "Алматы"; type = "keyboard"; name = "Razer Huntsman V3 Pro" },
         @{ loc = "Алматы"; type = "headset";  name = "SteelSeries Arctis Nova Pro" },
         @{ loc = "Алматы"; type = "monitor";  name = "Samsung Odyssey G7 27`"" }
       )
    },
    @{ name = "ABIX Esports Pack"; category = "gaming"; discount = 10
       pc = "ABIX Phantom RTX 5080"
       desc = "Сбалансированная сборка для турнирной игры и стримов."
       img  = "abix_rtx5080.webp"
       items = @(
         @{ loc = "Есиль"; type = "mouse";    name = "Razer DeathAdder V3" },
         @{ loc = "Есиль"; type = "keyboard"; name = "Keychron K2 Pro" },
         @{ loc = "Есиль"; type = "headset";  name = "HyperX Cloud III" },
         @{ loc = "Есиль"; type = "monitor";  name = "LG UltraGear 27GP850 27`"" }
       )
    },
    @{ name = "Streamer Studio Pro"; category = "streaming"; discount = 12
       pc = "Thermaltake Core Liquid"
       desc = "Максимум мощности для одновременной записи, кодирования и игры."
       img  = "thermaltake_core.webp"
       items = @(
         @{ loc = "Есиль"; type = "mouse";    name = "Razer DeathAdder V3" },
         @{ loc = "Есиль"; type = "keyboard"; name = "Keychron K2 Pro" },
         @{ loc = "Есиль"; type = "headset";  name = "HyperX Cloud III" },
         @{ loc = "Есиль"; type = "monitor";  name = "LG UltraGear 27GP850 27`"" }
       )
    },
    @{ name = "Budget Esports"; category = "gaming"; discount = 8
       pc = "ABIX Storm RTX 5060"
       desc = "Бюджетный сетап для CS2, Valorant и Dota 2 в 1080p 144 Гц."
       img  = "abix_rtx5060_aio.webp"
       items = @(
         @{ loc = "Сарыарка"; type = "mouse";    name = "Logitech G502 HERO" },
         @{ loc = "Сарыарка"; type = "keyboard"; name = "Logitech G413 SE" },
         @{ loc = "Сарыарка"; type = "headset";  name = "Razer Kraken X" },
         @{ loc = "Сарыарка"; type = "monitor";  name = "AOC 24G2 24`"" }
       )
    }
)

Write-Host "→ Creating $($Setups.Count) setups" -ForegroundColor Cyan
foreach ($s in $Setups) {
    $pcInfo = $createdPCs[$s.pc]
    if (-not $pcInfo) { Write-Warning "PC '$($s.pc)' not found"; continue }
    $perifIds = @()
    foreach ($it in $s.items) {
        $id = Get-Perif $it.loc $it.type $it.name
        if (-not $id) { Write-Warning "Peripheral '$($it.name)' not found"; continue }
        $perifIds += $id
    }
    $resp = Invoke-Api POST "/catalog/setups" @{
        name             = $s.name
        category         = $s.category
        pc_id            = $pcInfo.id
        peripheral_ids   = $perifIds
        discount_percent = $s.discount
        description      = $s.desc
        images           = @("/images/pcs/$($s.img)")
    }
    Write-Host "  + $($s.name) → $($resp.id)" -ForegroundColor Gray
}

Write-Host "`n✓ Done. Open http://localhost:3000 to see the catalog." -ForegroundColor Green
