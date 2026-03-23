function kClosest(points, k) {
    points.sort((a, b) => {
        const da = a[0] * a[0] + a[1] * a[1];
        const db = b[0] * b[0] + b[1] * b[1];
        if (da !== db) return da - db;
        if (a[0] !== b[0]) return a[0] - b[0];
        return a[1] - b[1];
    });
    return points.slice(0, k);
}

module.exports = { kClosest };
