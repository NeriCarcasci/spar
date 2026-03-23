use std::collections::BinaryHeap;
use std::cmp::Reverse;

pub struct MedianFinder {
    lo: BinaryHeap<i32>,
    hi: BinaryHeap<Reverse<i32>>,
}

impl MedianFinder {
    pub fn new() -> Self {
        MedianFinder {
            lo: BinaryHeap::new(),
            hi: BinaryHeap::new(),
        }
    }

    pub fn add_num(&mut self, num: i32) {
        self.lo.push(num);
        self.hi.push(Reverse(self.lo.pop().unwrap()));
        if self.hi.len() > self.lo.len() {
            self.lo.push(self.hi.pop().unwrap().0);
        }
    }

    pub fn find_median(&self) -> f64 {
        if self.lo.len() > self.hi.len() {
            *self.lo.peek().unwrap() as f64
        } else {
            (*self.lo.peek().unwrap() as f64 + self.hi.peek().unwrap().0 as f64) / 2.0
        }
    }
}
