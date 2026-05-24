use chrono::{Datelike, NaiveDate};

#[derive(Debug, Clone)]
pub struct Holiday {
    pub name: String,
    pub date: NaiveDate,
}

pub fn get_indonesian_holidays(year: i32) -> Vec<Holiday> {
    vec![
        Holiday { name: "Tahun Baru".into(), date: NaiveDate::from_ymd_opt(year, 1, 1).unwrap() },
        Holiday { name: "Isra Mi'raj".into(), date: NaiveDate::from_ymd_opt(year, 2, 8).unwrap() },
        Holiday { name: "Tahun Baru Imlek".into(), date: NaiveDate::from_ymd_opt(year, 1, 29).unwrap() },
        Holiday { name: "Hari Raya Nyepi".into(), date: NaiveDate::from_ymd_opt(year, 3, 29).unwrap() },
        Holiday { name: "Wafat Isa Almasih".into(), date: NaiveDate::from_ymd_opt(year, 4, 18).unwrap() },
        Holiday { name: "Hari Buruh".into(), date: NaiveDate::from_ymd_opt(year, 5, 1).unwrap() },
        Holiday { name: "Hari Raya Waisak".into(), date: NaiveDate::from_ymd_opt(year, 5, 12).unwrap() },
        Holiday { name: "Kenaikan Isa Almasih".into(), date: NaiveDate::from_ymd_opt(year, 5, 29).unwrap() },
        Holiday { name: "Hari Lahir Pancasila".into(), date: NaiveDate::from_ymd_opt(year, 6, 1).unwrap() },
        Holiday { name: "Idul Adha".into(), date: NaiveDate::from_ymd_opt(year, 6, 7).unwrap() },
        Holiday { name: "Tahun Baru Islam".into(), date: NaiveDate::from_ymd_opt(year, 6, 27).unwrap() },
        Holiday { name: "Hari Kemerdekaan".into(), date: NaiveDate::from_ymd_opt(year, 8, 17).unwrap() },
        Holiday { name: "Maulid Nabi".into(), date: NaiveDate::from_ymd_opt(year, 9, 5).unwrap() },
        Holiday { name: "Hari Natal".into(), date: NaiveDate::from_ymd_opt(year, 12, 25).unwrap() },
    ]
}

pub fn is_ramadan(date: NaiveDate) -> bool {
    (date >= NaiveDate::from_ymd_opt(2026, 2, 18).unwrap() && date <= NaiveDate::from_ymd_opt(2026, 3, 19).unwrap())
    || (date >= NaiveDate::from_ymd_opt(2027, 2, 8).unwrap() && date <= NaiveDate::from_ymd_opt(2027, 3, 9).unwrap())
}

pub fn get_sahoor_berbuka_times() -> (String, String) {
    // Approximate WIB times - in production, use prayer time API
    ("03:50 WIB".into(), "17:57 WIB".into())
}
