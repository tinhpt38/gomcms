package upload

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadHandler nhận file upload, lưu vào thư mục "uploads" và trả về URL ảnh.
func UploadHandler(c *gin.Context) {
	// Lấy file từ FormData với key "file"
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy file upload"})
		return
	}

	uploadDir := "./uploads"
	// Tạo thư mục nếu chưa tồn tại
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo thư mục lưu file"})
			return
		}
	}

	// Đặt tên file mới để tránh trùng lặp
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	filePath := filepath.Join(uploadDir, filename)

	// Lưu file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lưu file"})
		return
	}

	// Giả sử bạn đã cấu hình route để phục vụ các file tĩnh từ thư mục "uploads"
	photoURL := fmt.Sprintf("http://yourdomain.com/uploads/%s", filename)

	c.JSON(http.StatusOK, gin.H{"url": photoURL})
}
