package utils

// GetFirstItemFromMap returns the first key-value pair from the given map.
// If the map is empty, it returns zero values for the key and value.
func GetFirstItemFromMap[K comparable, V any](sourceMap map[K]V) (K, V, bool) {
    for key, value := range sourceMap {
        // Return the first encountered key-value pair along with true to indicate a successful find.
        return key, value, true
    }
    // Return zero values for the key and value types if the map is empty, along with false to indicate no items found.
    var zeroK K
    var zeroV V
    return zeroK, zeroV, false
}