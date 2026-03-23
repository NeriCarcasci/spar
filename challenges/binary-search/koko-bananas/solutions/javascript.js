function minEatingSpeed(piles, h) {
    let left = 1, right = Math.max(...piles);
    while (left < right) {
        const mid = left + Math.floor((right - left) / 2);
        const hours = piles.reduce((sum, p) => sum + Math.ceil(p / mid), 0);
        if (hours <= h) right = mid;
        else left = mid + 1;
    }
    return left;
}
module.exports = { minEatingSpeed };
