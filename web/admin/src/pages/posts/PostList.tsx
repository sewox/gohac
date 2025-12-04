import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Plus, Edit, Trash2, FileText } from 'lucide-react'
import toast from 'react-hot-toast'
import { postsAPI } from '../../lib/api'
import ConfirmDialog from '../../components/ConfirmDialog'
import '../pages/PageList.css'

interface Post {
  id: string
  slug: string
  title: string
  excerpt: string
  status: 'draft' | 'published' | 'archived'
  published_at: string | null
  created_at: string
  updated_at: string
}

export default function PostList() {
  const navigate = useNavigate()
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [deleteDialog, setDeleteDialog] = useState<{ open: boolean; post: Post | null }>({
    open: false,
    post: null,
  })

  useEffect(() => {
    fetchPosts()
  }, [])

  const fetchPosts = async () => {
    try {
      setLoading(true)
      const response = await postsAPI.list()
      setPosts(response.data.data || [])
      setError(null)
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to load posts'
      setError(errorMsg)
      toast.error(errorMsg)
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (post: Post) => {
    if (!post.id) return

    const deletePromise = postsAPI.delete(post.id)

    toast.promise(deletePromise, {
      loading: 'Deleting post...',
      success: 'Post deleted successfully!',
      error: (err: any) => err.response?.data?.error || 'Failed to delete post',
    })

    try {
      await deletePromise
      setDeleteDialog({ open: false, post: null })
      fetchPosts() // Refresh the list
    } catch (err: any) {
      console.error('Failed to delete post:', err)
    }
  }

  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'published':
        return 'status-badge published'
      case 'draft':
        return 'status-badge draft'
      case 'archived':
        return 'status-badge archived'
      default:
        return 'status-badge'
    }
  }

  if (loading) {
    return (
      <div className="page-list">
        <div className="loading">Loading posts...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="page-list">
        <div className="error-message">{error}</div>
      </div>
    )
  }

  return (
    <div className="page-list">
      <div className="page-list-header">
        <div className="header-title">
          <h1>Posts</h1>
        </div>
        <Link to="/admin/posts/new" className="create-button">
          <Plus size={20} />
          <span>New Post</span>
        </Link>
      </div>

      {posts.length === 0 ? (
        <div className="empty-state">
          <p>No posts yet. Create your first post to get started.</p>
          <Link to="/admin/posts/new" className="create-button">
            <Plus size={20} />
            <span>New Post</span>
          </Link>
        </div>
      ) : (
        <div className="page-list-table-container">
          <table className="page-list-table">
            <thead>
              <tr>
                <th>Title</th>
                <th>Slug</th>
                <th>Status</th>
                <th>Published</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {posts.map((post) => (
                <tr key={post.id}>
                  <td>
                    <strong>{post.title}</strong>
                    {post.excerpt && (
                      <div style={{ fontSize: '12px', color: '#666', marginTop: '4px' }}>
                        {post.excerpt.substring(0, 60)}...
                      </div>
                    )}
                  </td>
                  <td>
                    <code style={{ fontSize: '12px' }}>{post.slug}</code>
                  </td>
                  <td>
                    <span className={getStatusBadgeClass(post.status)}>{post.status}</span>
                  </td>
                  <td>
                    {post.published_at
                      ? new Date(post.published_at).toLocaleDateString()
                      : '-'}
                  </td>
                  <td>{new Date(post.created_at).toLocaleDateString()}</td>
                  <td className="actions">
                    <button
                      onClick={() => navigate(`/admin/posts/${post.id}/edit`)}
                      className="action-button edit"
                      title="Edit Post"
                    >
                      <Edit size={18} />
                    </button>
                    <button
                      onClick={() => setDeleteDialog({ open: true, post })}
                      className="action-button delete"
                      title="Delete Post"
                    >
                      <Trash2 size={18} />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <ConfirmDialog
        isOpen={deleteDialog.open}
        title="Delete Post"
        message={`Are you sure you want to delete "${deleteDialog.post?.title}"? This action cannot be undone.`}
        variant="danger"
        onConfirm={() => deleteDialog.post && handleDelete(deleteDialog.post)}
        onCancel={() => setDeleteDialog({ open: false, post: null })}
      />
    </div>
  )
}

