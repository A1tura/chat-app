package array;

type Array [T comparable] []T

func (array *Array[T]) ArrayInclude(target T) (bool, int) {
    for i, v := range *array {
        if v == target {
            return true, i
        }
    }

    return false, -1
}
