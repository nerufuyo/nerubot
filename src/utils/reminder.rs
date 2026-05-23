use chrono::{Datelike, NaiveDate, Utc};

#[derive(Debug, Clone)]
pub struct Holiday {
    pub name: String,
    pub date: NaiveDate,
    pub emoji: String,
}

pub fn get_indonesian_holidays(year: i32) -> Vec<Holiday> {
    vec![
        Holiday { name: "Tahun Baru".into(), date: NaiveDate::from_ymd_opt(year, 1, 1).unwrap(), emoji: "🎉".into() },
        Holiday { name: "Isra Mi'raj".into(), date: NaiveDate::from_ymd_opt(year, 2, 8).unwrap(), emoji: "🕌".into() },
        Holiday { name: "Tahun Baru Imlek".into(), date: NaiveDate::from_ymd_opt(year, 1, 29).unwrap(), emoji: "🧧".into() },
        Holiday { name: "Hari Raya Nyepi".into(), date: NaiveDate::from_ymd_opt(year, 3, 29).unwrap(), emoji: "🕉️".into() },
        Holiday { name: "Wafat Isa Almasih".into(), date: NaiveDate::from_ymd_opt(year, 4, 18).unwrap(), emoji: "✝️".into() },
        Holiday { name: "Hari Buruh".into(), date: NaiveDate::from_ymd_opt(year, 5, 1).unwrap(), emoji: "⚒️".into() },
        Holiday { name: "Hari Raya Waisak".into(), date: NaiveDate::from_ymd_opt(year, 5, 12).unwrap(), emoji: "☸️".into() },
        Holiday { name: "Kenaikan Isa Almasih".into(), date: NaiveDate::from_ymd_opt(year, 5, 29).unwrap(), emoji: "⬆️".into() },
        Holiday { name: "Hari Lahir Pancasila".into(), date: NaiveDate::from_ymd_opt(year, 6, 1).unwrap(), emoji: "🇮🇩".into() },
        Holiday { name: "Idul Adha".into(), date: NaiveDate::from_ymd_opt(year, 6, 7).unwrap(), emoji: "🐐".into() },
        Holiday { name: "Tahun Baru Islam".into(), date: NaiveDate::from_ymd_opt(year, 6, 27).unwrap(), emoji: "☪️".into() },
        Holiday { name: "Hari Kemerdekaan".into(), date: NaiveDate::from_ymd_opt(year, 8, 17).unwrap(), emoji: "🇮🇩".into() },
        Holiday { name: "Maulid Nabi".into(), date: NaiveDate::from_ymd_opt(year, 9, 5).unwrap(), emoji: "🕌".into() },
        Holiday { name: "Hari Natal".into(), date: NaiveDate::from_ymd_opt(year, 12, 25).unwrap(), emoji: "🎄".into() },
    ]
}

pub fn is_ramadan(date: NaiveDate) -> bool {
    // Approximate Ramadan dates - in production, use Hijri calendar library
    let year = date.year();
    match year {
        2026 => date >= NaiveDate::from_ymd_opt(2026, 2, 18).unwrap()
             && date <= NaiveDate::from_ymd_opt(2026, 3, 19).unwrap(),
        2027 => date >= NaiveDate::from_ymd_opt(2027, 2, 8).unwrap()
             && date <= NaiveDate::from_ymd_opt(2027, 3, 9).unwrap(),
        _ => false,
    }
}

pub fn get_sahoor_berbuka_times() -> (String, String) {
    // Approximate WIB times - in production, use prayer time API
    ("03:50 WIB".into(), "17:57 WIB".into())
}
