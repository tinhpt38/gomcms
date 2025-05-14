// Search functionality
const searchQuery = ref('')

// Computed property for filtered history data
const filteredHistory = computed(() => {
    if (!searchQuery.value) {
        return luckyHistory.value
    }
    const query = searchQuery.value.trim().toLowerCase()
    return luckyHistory.value.filter(item => 
        item.luckyNumber.toString().includes(query)
    )
})

// Computed property for paginated history data
const paginatedHistory = computed(() => {
    const startIndex = (currentPage.value - 1) * pageSize.value
    const endIndex = startIndex + pageSize.value
    return filteredHistory.value.slice(startIndex, endIndex)
})

// Pagination event handlers
const handleSizeChange = (val) => {
    pageSize.value = val
    currentPage.value = 1 // Reset to first page when changing page size
}

const handleCurrentChange = (val) => {
    currentPage.value = val
}

// Handle search function
const handleSearch = () => {
    // Reset to first page when searching
    currentPage.value = 1
    // The actual filtering happens in the filteredHistory computed property
}

// Reset search function
const resetSearch = () => {
    searchQuery.value = ''
    currentPage.value = 1
}
