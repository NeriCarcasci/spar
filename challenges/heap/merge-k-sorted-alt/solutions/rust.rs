pub fn k_closest(mut points: Vec<Vec<i32>>, k: i32) -> Vec<Vec<i32>> {
    let k = k as usize;
    points.sort_by(|a, b| {
        let da = a[0] * a[0] + a[1] * a[1];
        let db = b[0] * b[0] + b[1] * b[1];
        da.cmp(&db).then(a[0].cmp(&b[0])).then(a[1].cmp(&b[1]))
    });
    points.truncate(k);
    points
}
