package util

import "golang.org/x/exp/constraints"

func MaxOf[T constraints.Ordered](vars ...T) T {
    max := vars[0]

    for _, i := range vars {
        if max < i {
            max = i
        }
    }

    return max
}

func MinOf[T constraints.Ordered](vars ...T) T {
	min := vars[0]

    for _, i := range vars {
        if min > i {
            min = i
        }
    }

    return min
}