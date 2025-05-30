package num

import "fmt"

// Matrix3 defines a 3x3 matrix of float64 values
type Matrix3 [3]Vector3

func NewMatrix3(x1, y1, z1, x2, y2, z2, x3, y3, z3 float64) Matrix3 {
	return Matrix3{
		{x1, y1, z1},
		{x2, y2, z2},
		{x3, y3, z3},
	}
}

// MultiplyXYZ takes x, y, z and creates 3D Vector3. Then Multiply m with the
// newly created Vector3. Returns the resulting vector
func (m Matrix3) MultiplyXYZ(x, y, z float64) Vector3 {
	return m.Multiply(NewVector3(x, y, z))
}

// Multiply applies the matrix to a 3D vector and returns the resulting vector
func (m Matrix3) Multiply(v Vector3) Vector3 {
	var result Vector3
	for i := range 3 {
		for j := range 3 {
			result[i] += m[i][j] * v[j]
		}
	}
	return result
}

// Transpose transposes the Matrix3
func (m Matrix3) Transpose() Matrix3 {
	var result Matrix3
	for i := range 3 {
		for j := range 3 {
			result[j][i] = m[i][j]
		}
	}
	return result
}

func (m Matrix3) Inverse() (Matrix3, bool) {
	a, b, c := m[0][0], m[0][1], m[0][2]
	d, e, f := m[1][0], m[1][1], m[1][2]
	g, h, i := m[2][0], m[2][1], m[2][2]

	// Compute the determinant
	det := a*(e*i-f*h) - b*(d*i-f*g) + c*(d*h-e*g)
	if det == 0 {
		return Matrix3{}, false // Matrix is not invertible
	}
	invDet := 1.0 / det

	var inv Matrix3
	inv[0][0] = (e*i - f*h) * invDet
	inv[0][1] = -(b*i - c*h) * invDet
	inv[0][2] = (b*f - c*e) * invDet
	inv[1][0] = -(d*i - f*g) * invDet
	inv[1][1] = (a*i - c*g) * invDet
	inv[1][2] = -(a*f - c*d) * invDet
	inv[2][0] = (d*h - e*g) * invDet
	inv[2][1] = -(a*h - b*g) * invDet
	inv[2][2] = (a*e - b*d) * invDet

	return inv, true
}

func (m Matrix3) String() string {
	return fmt.Sprintf("[\n\t%.10f,%.10f,%.10f,\n\t%.10f,%.10f,%.10f,\n\t%.10f,%.10f,%.10f,\n]",
		m[0][0], m[0][1], m[0][2],
		m[1][0], m[1][1], m[1][2],
		m[2][0], m[2][1], m[2][2],
	)
}

// Vector3 defines a 3D vector
type Vector3 [3]float64

// NewVector3 create new 3D vector: Vector3
func NewVector3(x, y, z float64) Vector3 {
	return Vector3{x, y, z}
}

func (v Vector3) MultiplyMatrix(m Matrix3) Vector3 {
	var result Vector3
	for j := range 3 {
		for i := range 3 {
			result[j] += v[i] * m[i][j]
		}
	}
	return result
}

func (v Vector3) MultiplyScalar(s float64) Vector3 {
	var result Vector3
	for i := range 3 {
		result[i] = v[i] * s
	}
	return result
}

func (v Vector3) Add(vec Vector3) Vector3 {
	var result Vector3
	for i := range 3 {
		result[i] = v[i] + vec[i]
	}
	return result
}

func (v Vector3) Values() (float64, float64, float64) {
	return v[0], v[1], v[2]
}
