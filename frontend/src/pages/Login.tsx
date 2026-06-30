import { useForm } from 'react-hook-form'
import { Button } from 'primereact/button'
import { Card } from 'primereact/card'
import { useNavigate } from 'react-router'
import { useAuth } from '../auth'
import { Input } from '../component/input/Input'

interface LoginForm {
  email: string
  password: string
}

export default function Login() {
  const navigate = useNavigate()
  const { login } = useAuth()

  const {
    control,
    handleSubmit,
    formState: { errors, isValid },
  } = useForm<LoginForm>({ mode: 'onChange' })

  function onSubmit(data: LoginForm) {
    console.log('Login attempt', data)
    login()
    navigate('/dashboard')
  }

  return (
    <div className="flex align-items-center justify-content-center min-h-screen">
      <Card title="Sign In" subTitle="Beholder" className="w-4">
        <form
          onSubmit={handleSubmit(onSubmit)}
          className="flex flex-column gap-4 p-4"
        >
          <Input control={control} name="email" label="Username" type="text"
                 rules={{ required: 'Email is required' }} errors={errors} />

          <Input control={control} name="password" label="Password" type="password"
                 rules={{ required: 'Password is required' }} errors={errors} />

          <Button
            type="submit"
            label="Sign In"
            icon="pi pi-sign-in"
            disabled={!isValid}
            className="w-full"
          />

          <Button
            type="button"
            label="Don't have an account? Sign up"
            link
            className="w-full"
            onClick={() => navigate('/signup')}
          />
        </form>
      </Card>
    </div>
  )
}
