package auth

type LoginDTO struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type RegisterDTO struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,password"`
    Name     string `json:"name" binding:"required,gt=0"`        
}