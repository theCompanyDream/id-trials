import React, { useState, useMemo } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { useUserStore } from '../components';

const UserDetail = () => {
  const { id } = useParams(); // If there's an ID, we're editing an existing user.
  const navigate = useNavigate(); // For programmatic navigation
  const userId = useUserStore((state) => state.userId)

  // Initialize form data with empty strings.
  const [formData, setFormData] = useState({
    hash: '',
    user_name: '',
    first_name: '',
    last_name: '',
    email: '',
    department: '',
  });

  // State for handling errors
  const [errors, setErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  // If editing, fetch the existing user data.
  useMemo(() => {
    if (id) {
      fetch(`/api/${userId}/${id}`)
        .then((res) => {
          if (!res.ok) throw new Error('Error fetching user data');
          return res.json();
        })
        .then((data) => {
          setFormData({
            hash: data.hash || '',
            user_name: data.user_name || '',
            first_name: data.first_name || '',
            last_name: data.last_name || '',
            email: data.email || '',
            user_status: data.user_status || '',
            department: data.department || '',
          });
        })
        .catch((err) => {
          console.error('Error:', err);
          setErrors({ general: 'Failed to load user data' });
        });
    }
  }, [id]);

  const handleSelect = (e) => {
    updateStore({
      userId: e.target.value
    })
  }

  // Handle input changes.
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));

    // Clear field-specific error when user starts typing
    if (errors[name]) {
      setErrors((prev) => {
        const newErrors = { ...prev };
        delete newErrors[name];
        return newErrors;
      });
    }
  };

  // Handle form submission.
  const handleSubmit = (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setErrors({}); // Clear previous errors

    const url = id ? `/api/${userId}/${id}` : `/api/${userId}`;
    const method = id ? 'PUT' : 'POST';

    fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(formData),
    })
      .then(async (res) => {
        const data = await res.json();

        if (!res.ok) {
          // Handle validation errors from Echo (status 422)
          if (res.status === 422) {
            // Echo returns validation errors directly as an object map
            setErrors(data);
          } else {
            // Handle other types of errors
            setErrors({ general: data.message || 'Error saving user' });
          }
          return;
        }
        setErrors({});

        // Success - navigate to the user detail page
        // Use the ID field from the UserDTO response
        const userId = data.id || id; // data.ID is the primary key from your UserDTO
        navigate(`/detail/${userId}`);

        return data;
      })
      .catch((err) => {
        console.error('Error:', err);
        // Only set general error if no specific errors were already set
        if (!errors.general && Object.keys(errors).length === 0) {
          setErrors({ general: 'An unexpected error occurred' });
        }
      })
      .finally(() => {
        setIsSubmitting(false);
      });
  };

  return (
    <main className="max-w-xl mx-auto p-4">
      <h2 className="text-2xl font-bold mb-4">
        {id ? 'Edit User' : 'Create User'}
      </h2>

      {/* Display general errors */}
      {errors && Object.keys(errors).length > 0 && (
        <section className="bg-red-50 border-l-4 border-red-400 p-4 mb-6">
          <div className="text-sm text-red-700">
            <p className="font-medium text-red-800 mb-2">Please correct the following errors:</p>
            <ul className="list-disc list-inside space-y-1">
              {Object.entries(errors).map(([field, message]) => (
                <li key={field}>
                  <span className="font-medium">{field}:</span> {message}
                </li>
              ))}
            </ul>
          </div>
        </section>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
      {!id && (<div>
          <label className="block mb-1 font-semibold">Id</label>
          <select
            onChange={handleSelect}
            value={userId}
            className="px-4 py-2 border rounded-lg bg-white w-full font-medium focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition hover:border-blue-400 cursor-pointer"
          >
            <option value="uuid4">UUID</option>
            <option value="cuidId">CUID</option>
            <option value="snowId">Snowflake</option>
            <option value="ksuidId">KSUID</option>
            <option value="ulidId">ULID</option>
            <option value="nanoId">NanoID</option>
          </select>
          {errors.user_name && (
            <p className="text-red-600 text-sm mt-1">{errors.user_name}</p>
          )}
        </div>)}

        <div>
          <label className="block mb-1 font-semibold">Username</label>
          <input
            type="text"
            name="user_name"
            value={formData.user_name}
            onChange={handleChange}
            className={`w-full border rounded p-2 ${
              errors.user_name ? 'border-red-400' : 'border-gray-300'
            }`}
            required
          />
          {errors.user_name && (
            <p className="text-red-600 text-sm mt-1">{errors.user_name}</p>
          )}
        </div>

        <div>
          <label className="block mb-1 font-semibold">First Name</label>
          <input
            type="text"
            name="first_name"
            value={formData.first_name}
            onChange={handleChange}
            className={`w-full border rounded p-2 ${
              errors.first_name ? 'border-red-400' : 'border-gray-300'
            }`}
            required
          />
          {errors.first_name && (
            <p className="text-red-600 text-sm mt-1">{errors.first_name}</p>
          )}
        </div>

        <div>
          <label className="block mb-1 font-semibold">Last Name</label>
          <input
            type="text"
            name="last_name"
            value={formData.last_name}
            onChange={handleChange}
            className={`w-full border rounded p-2 ${
              errors.last_name ? 'border-red-400' : 'border-gray-300'
            }`}
            required
          />
          {errors.last_name && (
            <p className="text-red-600 text-sm mt-1">{errors.last_name}</p>
          )}
        </div>

        <div>
          <label className="block mb-1 font-semibold">Email</label>
          <input
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            className={`w-full border rounded p-2 ${
              errors.email ? 'border-red-400' : 'border-gray-300'
            }`}
            required
          />
          {errors.email && (
            <p className="text-red-600 text-sm mt-1">{errors.email}</p>
          )}
        </div>

        <div>
          <label className="block mb-1 font-semibold">Department</label>
          <select
            name="department"
            value={formData.department}
            onChange={handleChange}
            className={`w-full border rounded p-2 ${
              errors.department ? 'border-red-400' : 'border-gray-300'
            }`}
          >
            <option value="">Select a Department</option>
            <option value="accounting">Accounting</option>
            <option value="it">Information Technology</option>
            <option value="hr">Human Resources</option>
            <option value="marketing">Marketing</option>
            <option value="sales">Sales</option>
            <option value="engineering">Engineering</option>
            <option value="operations">Operations</option>
          </select>
          {errors.department && (
            <p className="text-red-600 text-sm mt-1">{errors.department}</p>
          )}
        </div>

        <div className="flex gap-4">
          <button
            type="submit"
            disabled={isSubmitting}
            className="bg-blue-500 hover:bg-blue-600 disabled:bg-blue-300 text-white px-4 py-2 rounded"
          >
            {isSubmitting ? 'Saving...' : (id ? 'Update User' : 'Create User')}
          </button>

          <Link
            className="bg-gray-100 hover:bg-gray-200 text-black px-4 py-2 rounded border"
            to="/"
          >
            Back
          </Link>
        </div>
      </form>
    </main>
  );
};

export default UserDetail;