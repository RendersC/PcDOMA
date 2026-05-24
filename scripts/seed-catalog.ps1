# PCDoma catalog seed script — new collection.
# Logs in as admin, wipes ALL existing PCs / peripherals / setups, and
# re-creates a curated catalog from the new photo collection:
#   12 PCs (setup_01..setup_12) + 11 ready-made setups (setup_13..setup_23).
# Peripherals are intentionally left empty (catalog reset to the new set).
#
# Usage:   powershell -ExecutionPolicy Bypass -File scripts\seed-catalog.ps1
# Assumes: docker compose is up; gateway on http://localhost:8090

$ErrorActionPreference = "Stop"
$Gateway              = "http://localhost:8090/api/v1"
$AdminEmail           = "test@pcdoma.kz"
$AdminPassword        = "password123"
$MongoContainer       = "pcdoma-mongo"
$MongoUser            = "pcdoma"
$MongoPassword        = "secret"

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
try { $existingSetups = (Invoke-Api GET "/catalog/setups?limit=200" $null).data } catch { $existingSetups = @() }
foreach ($s in $existingSetups) {
    try { Invoke-Api DELETE "/catalog/setups/$($s.id)" $null | Out-Null } catch {}
}
try { $existingPeripherals = (Invoke-Api GET "/catalog/peripherals?limit=200" $null).data } catch { $existingPeripherals = @() }
foreach ($p in $existingPeripherals) {
    try { Invoke-Api DELETE "/catalog/peripherals/$($p.id)" $null | Out-Null } catch {}
}
try { $existingPCs = (Invoke-Api GET "/catalog/pcs?limit=200" $null).data } catch { $existingPCs = @() }
foreach ($pc in $existingPCs) {
    try { Invoke-Api DELETE "/catalog/pcs/$($pc.id)" $null | Out-Null } catch {}
}
Write-Host "  ✓ wiped $($existingPCs.Count) PCs, $($existingPeripherals.Count) peripherals, $($existingSetups.Count) setups" -ForegroundColor Green

# ── 3. Locations (3 districts of Astana) ────────────────────────────
$LocEsil     = @{ district = "Есиль";    address = "ул. Кунаева 14";    lat = 51.1284; lng = 71.4310 }
$LocAlmaty   = @{ district = "Алматы";   address = "пр. Республики 52"; lat = 51.1605; lng = 71.4704 }
$LocSaryarka = @{ district = "Сарыарка"; address = "ул. Иманова 19";    lat = 51.1773; lng = 71.4491 }

# ── 4. PCs (12 builds, setup_01..setup_12) ──────────────────────────
$PCs = @(
    @{ name = "Nebula RGB X";    img = "setup_01.jpg"; loc = $LocEsil;
       cpu = "Intel Core i9-14900K";  gpu = "NVIDIA RTX 4080 SUPER 16GB";    ram = "64GB DDR5 6400 MHz"; storage = "2TB NVMe Samsung 990 Pro";
       hour = 5600; day = 16000 },
    @{ name = "Phantom White";   img = "setup_02.jpg"; loc = $LocAlmaty;
       cpu = "Intel Core i7-14700K";  gpu = "NVIDIA RTX 4070 Ti SUPER 16GB"; ram = "32GB DDR5 6000 MHz"; storage = "1TB NVMe WD Black SN850";
       hour = 4600; day = 13500 },
    @{ name = "Crimson Strike";  img = "setup_03.jpg"; loc = $LocSaryarka;
       cpu = "AMD Ryzen 7 7800X3D";   gpu = "NVIDIA RTX 4070 SUPER 12GB";    ram = "32GB DDR5 6000 MHz"; storage = "1TB NVMe Kingston KC3000";
       hour = 4000; day = 11800 },
    @{ name = "Frost Compact";   img = "setup_04.jpg"; loc = $LocEsil;
       cpu = "Intel Core i5-13600KF"; gpu = "NVIDIA RTX 4070 12GB";          ram = "32GB DDR5 5600 MHz"; storage = "1TB NVMe Samsung 980 Pro";
       hour = 3600; day = 10500 },
    @{ name = "Aurora Mini";     img = "setup_05.jpg"; loc = $LocSaryarka;
       cpu = "AMD Ryzen 5 7600";      gpu = "NVIDIA RTX 4060 8GB";           ram = "16GB DDR5 5200 MHz"; storage = "512GB NVMe Kingston NV2";
       hour = 2200; day = 6500 },
    @{ name = "SimRacer Pro";    img = "setup_06.jpg"; loc = $LocAlmaty;
       cpu = "Intel Core i7-13700K";  gpu = "NVIDIA RTX 4070 SUPER 12GB";    ram = "32GB DDR5 6000 MHz"; storage = "1TB NVMe Samsung 990 Pro";
       hour = 4200; day = 12500 },
    @{ name = "Onyx Tower";      img = "setup_07.jpg"; loc = $LocEsil;
       cpu = "Intel Core i9-13900K";  gpu = "NVIDIA RTX 4080 16GB";          ram = "64GB DDR5 6400 MHz"; storage = "2TB NVMe Samsung 990 Pro";
       hour = 5400; day = 15500 },
    @{ name = "Violet Haze";     img = "setup_08.jpg"; loc = $LocAlmaty;
       cpu = "AMD Ryzen 7 7700X";     gpu = "NVIDIA RTX 4070 Ti 12GB";       ram = "32GB DDR5 6000 MHz"; storage = "1TB NVMe Samsung 980 Pro";
       hour = 3800; day = 11200 },
    @{ name = "Carbon Edge";     img = "setup_09.jpg"; loc = $LocSaryarka;
       cpu = "Intel Core i5-14400F";  gpu = "NVIDIA RTX 4060 Ti 8GB";        ram = "16GB DDR5 5600 MHz"; storage = "1TB NVMe Crucial P5 Plus";
       hour = 2800; day = 8200 },
    @{ name = "Titan Liquid";    img = "setup_10.jpg"; loc = $LocEsil;
       cpu = "Intel Core i9-14900KS"; gpu = "NVIDIA RTX 4090 24GB";          ram = "128GB DDR5 6400 MHz"; storage = "4TB NVMe Samsung 990 Pro + 8TB HDD";
       hour = 8000; day = 22000 },
    @{ name = "Glacier Pro";     img = "setup_11.jpg"; loc = $LocAlmaty;
       cpu = "AMD Ryzen 9 7950X3D";   gpu = "NVIDIA RTX 4090 24GB";          ram = "64GB DDR5 6400 MHz"; storage = "4TB NVMe Samsung 990 Pro";
       hour = 7000; day = 19500 },
    @{ name = "Stealth Black";   img = "setup_12.jpg"; loc = $LocSaryarka;
       cpu = "Intel Core i7-14700F";  gpu = "NVIDIA RTX 4070 12GB";          ram = "32GB DDR5 5600 MHz"; storage = "1TB NVMe Samsung 980 Pro";
       hour = 3400; day = 10000 }
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
    $createdPCs[$p.name] = $resp.id
    Write-Host "  + $($p.name) → $($resp.id)" -ForegroundColor Gray
}

# ── 4b. Peripherals (22: 8 headsets, 7 mice, 7 keyboards, with photos) ─
$Peripherals = @(
    # ── Наушники (headphones) ──
    @{ name = "Razer BlackShark V2 X"; type = "headphones"; loc = "Есиль";   img = "headset_01.jpg";
       desc = "Лёгкая киберспортивная гарнитура с объёмным звуком 7.1.";
       specs = @{ driver = "50mm"; surround = "7.1"; weight = "240g" }; hour = 200; day = 900 },
    @{ name = "HyperX Cloud III";      type = "headphones"; loc = "Алматы";  img = "headset_02.jpg";
       desc = "Комфортная гарнитура с DTS Headphone:X и чётким микрофоном.";
       specs = @{ driver = "53mm"; mic = "DTS"; weight = "320g" }; hour = 220; day = 1000 },
    @{ name = "SteelSeries Arctis Nova 7"; type = "headphones"; loc = "Сарыарка"; img = "headset_03.jpg";
       desc = "Беспроводная гарнитура Hi-Fi с двойной беспроводной связью.";
       specs = @{ driver = "40mm"; wireless = "2.4G + BT"; battery = "38h" }; hour = 320; day = 1500 },
    @{ name = "Aimzone Cat Ear RGB";   type = "headphones"; loc = "Есиль";   img = "headset_04.jpg";
       desc = "Стильная гарнитура с кошачьими ушками и RGB-подсветкой.";
       specs = @{ driver = "50mm"; rgb = "yes"; mic = "съёмный" }; hour = 250; day = 1200 },
    @{ name = "Stereo Gaming Pro";     type = "headphones"; loc = "Сарыарка"; img = "headset_05.jpg";
       desc = "Бюджетная стереогарнитура для долгих игровых сессий.";
       specs = @{ driver = "40mm"; surround = "Stereo"; weight = "280g" }; hour = 150; day = 700 },
    @{ name = "Logitech G PRO X";      type = "headphones"; loc = "Алматы";  img = "headset_06.jpg";
       desc = "Турнирная гарнитура с микрофоном Blue VO!CE.";
       specs = @{ driver = "50mm"; mic = "Blue VO!CE"; weight = "320g" }; hour = 280; day = 1400 },
    @{ name = "HyperX Cloud II";       type = "headphones"; loc = "Есиль";   img = "headset_07.jpg";
       desc = "Легендарная гарнитура с виртуальным звуком 7.1.";
       specs = @{ driver = "53mm"; surround = "7.1"; weight = "320g" }; hour = 210; day = 950 },
    @{ name = "Razer Kraken V2 7.1";   type = "headphones"; loc = "Алматы";  img = "headset_08.jpg";
       desc = "Гарнитура с охлаждающими амбушюрами и THX Spatial Audio.";
       specs = @{ driver = "50mm"; surround = "THX 7.1"; weight = "322g" }; hour = 230; day = 1100 },

    # ── Мышки (mouse) ──
    @{ name = "Razer DeathAdder V3";   type = "mouse"; loc = "Есиль";   img = "mouse_01.jpg";
       desc = "Эргономичная киберспортивная мышь, сенсор Focus Pro 30K.";
       specs = @{ sensor = "Focus Pro 30K"; dpi = "30000"; weight = "59g" }; hour = 150; day = 700 },
    @{ name = "Logitech G Pro X Superlight 2"; type = "mouse"; loc = "Алматы"; img = "mouse_02.jpg";
       desc = "Беспроводная сверхлёгкая мышь, сенсор HERO 2.";
       specs = @{ sensor = "HERO 2"; dpi = "32000"; weight = "60g"; wireless = "LIGHTSPEED" }; hour = 200; day = 950 },
    @{ name = "Glorious Model O";      type = "mouse"; loc = "Сарыарка"; img = "mouse_03.jpg";
       desc = "Ультралёгкая мышь с сотовым корпусом, 67 г.";
       specs = @{ sensor = "BAMF"; dpi = "19000"; weight = "67g" }; hour = 130; day = 650 },
    @{ name = "Pulsar X2 Ergo";        type = "mouse"; loc = "Есиль";   img = "mouse_04.jpg";
       desc = "Лёгкая эргономичная мышь с оптическим сенсором.";
       specs = @{ sensor = "PAW3395"; dpi = "26000"; weight = "54g" }; hour = 140; day = 680 },
    @{ name = "Endgame Gear XM2";      type = "mouse"; loc = "Алматы";  img = "mouse_05.jpg";
       desc = "Точная и резкая мышь для соревновательной игры.";
       specs = @{ sensor = "PAW3395"; dpi = "26000"; weight = "63g" }; hour = 140; day = 680 },
    @{ name = "MuMCHOSE AX5 Pro Max";  type = "mouse"; loc = "Сарыарка"; img = "mouse_06.jpg";
       desc = "Беспроводная мышь с опросом 8 кГц и сенсором PAW3950.";
       specs = @{ sensor = "PAW3950"; dpi = "42000"; polling = "8000Hz"; weight = "49g" }; hour = 180; day = 850 },
    @{ name = "Logitech G502 LIGHTSPEED"; type = "mouse"; loc = "Есиль"; img = "mouse_07.jpg";
       desc = "Классическая мышь с настраиваемыми грузами, беспроводная.";
       specs = @{ sensor = "HERO 25K"; dpi = "25600"; weight = "114g"; wireless = "LIGHTSPEED" }; hour = 160; day = 780 },

    # ── Клавиатуры (keyboard) ──
    @{ name = "Aluminium Gasket 75%";  type = "keyboard"; loc = "Есиль";   img = "keyboard_01.jpg";
       desc = "Премиальная gasket-клавиатура в алюминиевом корпусе.";
       specs = @{ layout = "75%"; mount = "Gasket"; switch = "Linear"; backlight = "RGB" }; hour = 300; day = 1500 },
    @{ name = "Keychron K2 Pro";       type = "keyboard"; loc = "Алматы";  img = "keyboard_02.jpg";
       desc = "75% механика с hot-swap и Gateron G Pro.";
       specs = @{ layout = "75%"; switch = "Gateron G Pro"; hotswap = "yes"; backlight = "RGB" }; hour = 250; day = 1200 },
    @{ name = "Royal Kludge RK84";     type = "keyboard"; loc = "Сарыарка"; img = "keyboard_03.jpg";
       desc = "Беспроводная 75% клавиатура с тройным подключением.";
       specs = @{ layout = "75%"; switch = "RK Brown"; wireless = "BT+2.4G+USB" }; hour = 200; day = 950 },
    @{ name = "Samurai PBT Custom";    type = "keyboard"; loc = "Есиль";   img = "keyboard_04.jpg";
       desc = "Кастом с тематическими PBT-кейкапами в стиле самурая.";
       specs = @{ layout = "65%"; keycaps = "PBT Dye-sub"; switch = "Tactile" }; hour = 280; day = 1350 },
    @{ name = "Mercury K1 Gradient White"; type = "keyboard"; loc = "Алматы"; img = "keyboard_05.jpg";
       desc = "75% игровая клавиатура с градиентным белым дизайном.";
       specs = @{ layout = "75%"; switch = "Red Linear"; backlight = "RGB" }; hour = 270; day = 1300 },
    @{ name = "REZE Anime PBT Custom"; type = "keyboard"; loc = "Сарыарка"; img = "keyboard_06.jpg";
       desc = "Кастомная сборка с аниме-кейкапами PBT.";
       specs = @{ layout = "65%"; keycaps = "PBT"; switch = "Linear" }; hour = 260; day = 1250 },
    @{ name = "Teclado Mecânico Pro";  type = "keyboard"; loc = "Есиль";   img = "keyboard_07.jpg";
       desc = "Полноразмерная механическая клавиатура с RGB.";
       specs = @{ layout = "Full"; switch = "Blue"; backlight = "RGB" }; hour = 180; day = 850 }
)

Write-Host "→ Creating $($Peripherals.Count) peripherals" -ForegroundColor Cyan
foreach ($p in $Peripherals) {
    $resp = Invoke-Api POST "/catalog/peripherals" @{
        name           = $p.name
        type           = $p.type
        description    = $p.desc
        specs          = $p.specs
        price_per_hour = $p.hour
        price_per_day  = $p.day
        location_id    = $p.loc
        images         = @("/images/peripherals/$($p.img)")
    }
    Write-Host "  + $($p.name) [$($p.type)] → $($resp.id)" -ForegroundColor Gray
}

# ── 5. Setups (11, setup_13..setup_23) ──────────────────────────────
# Each setup references a real PC; peripheral list is empty in this reset.
$Setups = @(
    @{ name = "Gojo Aesthetic Setup"; category = "gaming";    discount = 10; pc = "Titan Liquid";  img = "setup_13.jpg";
       desc = "Футуристичный сетап в стиле Годжо: топовое железо и чистый стол для максимального погружения." },
    @{ name = "Matte Creator Studio"; category = "design";    discount = 10; pc = "Onyx Tower";    img = "setup_14.jpg";
       desc = "Матово-чёрная рабочая станция для монтажа, 3D и креативных задач." },
    @{ name = "White Pearl Gaming";   category = "gaming";    discount = 12; pc = "Phantom White"; img = "setup_15.jpg";
       desc = "Эстетичная белая сборка — мощно и стильно для любой игры." },
    @{ name = "Minimal Clean Desk";   category = "work";      discount = 8;  pc = "Frost Compact"; img = "setup_16.jpg";
       desc = "Минималистичный лёгкий сетап для работы и учёбы без лишнего." },
    @{ name = "Pro Performance Bundle"; category = "gaming";  discount = 10; pc = "Nebula RGB X";  img = "setup_17.jpg";
       desc = "Сбалансированный RGB-сетап для AAA-проектов и киберспорта." },
    @{ name = "Magic Tech Build";     category = "gaming";    discount = 12; pc = "Glacier Pro";   img = "setup_18.jpg";
       desc = "Флагманская сборка с RTX 4090 для 4K-гейминга на ультрах." },
    @{ name = "Insta Streamer Kit";   category = "streaming"; discount = 12; pc = "Crimson Strike"; img = "setup_19.jpg";
       desc = "Готовая станция для стримов: запись, кодирование и игра одновременно." },
    @{ name = "Keyboard Paradise";    category = "gaming";    discount = 8;  pc = "Violet Haze";   img = "setup_20.jpg";
       desc = "Сетап для тех, кто ценит механику и отзывчивость в каждом клике." },
    @{ name = "Mercado Gamer";        category = "gaming";    discount = 10; pc = "SimRacer Pro";  img = "setup_21.jpg";
       desc = "Универсальный игровой сетап для симрейсинга и сетевых баталий." },
    @{ name = "Ultimate Link Setup";  category = "gaming";    discount = 10; pc = "Stealth Black"; img = "setup_22.jpg";
       desc = "Тёмный сетап с акцентной подсветкой для комфортных долгих сессий." },
    @{ name = "Full Gear Station";    category = "streaming"; discount = 12; pc = "Aurora Mini";   img = "setup_23.jpg";
       desc = "Компактная станция со всем необходимым для контента и игры." }
)

Write-Host "→ Creating $($Setups.Count) setups" -ForegroundColor Cyan
foreach ($s in $Setups) {
    $pcId = $createdPCs[$s.pc]
    if (-not $pcId) { Write-Warning "PC '$($s.pc)' not found"; continue }
    $resp = Invoke-Api POST "/catalog/setups" @{
        name             = $s.name
        category         = $s.category
        pc_id            = $pcId
        peripheral_ids   = @()
        discount_percent = $s.discount
        description      = $s.desc
        images           = @("/images/pcs/$($s.img)")
    }
    Write-Host "  + $($s.name) → $($resp.id)" -ForegroundColor Gray
}

# ── 6. Double setup prices ──────────────────────────────────────────
# Setup price is auto-derived from the PC price on creation, so to make
# setups twice as expensive (without touching PC catalog prices) we
# multiply the stored totals directly in MongoDB. Keeps discount badges.
Write-Host "→ Doubling setup prices (x2)" -ForegroundColor Cyan
$mongoUri  = "mongodb://${MongoUser}:${MongoPassword}@localhost:27017/pcrentalcatalog?authSource=admin"
# Uses the $mul update operator (no field-path strings / no embedded quotes),
# so PowerShell 5.1 passes the --eval argument to docker intact.
$mongoEval = 'db.setups.updateMany({},{$mul:{total_price_per_hour:2,total_price_per_day:2}}).modifiedCount'
$modified  = docker exec $MongoContainer mongosh $mongoUri --quiet --eval $mongoEval
Write-Host "  x2 applied to $modified setups" -ForegroundColor Gray

Write-Host "`n✓ Done. Open http://localhost:3000 to see the new collection." -ForegroundColor Green
