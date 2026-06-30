import { Button } from 'primereact/button'
import { Card } from 'primereact/card'
import { Avatar } from 'primereact/avatar'
import { Chip } from 'primereact/chip'
import { useNavigate } from 'react-router'
import { useAuth } from '../auth'

export default function Dashboard() {
  const navigate = useNavigate()
  const { logout } = useAuth()

  function handleLogout() {
    logout()
    navigate('/')
  }

  return (
    <div className="min-h-screen flex flex-column">
      <header className="flex align-items-center justify-content-between px-5 py-3 surface-0 shadow-2">
        <div className="flex align-items-center gap-3">
          <i className="pi pi-shield text-3xl text-primary" />
          <span className="text-xl font-medium">Beholder</span>
        </div>
        <div className="flex align-items-center gap-3">
          <Chip
            label="admin@beholder.app"
            template={<Avatar icon="pi pi-user" shape="circle" className="mr-2" />}
          />
          <Button
            icon="pi pi-sign-out"
            severity="secondary"
            text
            onClick={handleLogout}
            tooltip="Sign out"
          />
        </div>
      </header>

      <main className="flex-1 p-5">
        <div className="grid">
          <div className="col-12 md:col-6 lg:col-3">
            <Card className="text-center">
              <i className="pi pi-users text-4xl text-primary mb-3 block" />
              <h2 className="text-2xl font-medium m-0">1,234</h2>
              <p className="text-600 mt-2">Active Users</p>
            </Card>
          </div>
          <div className="col-12 md:col-6 lg:col-3">
            <Card className="text-center">
              <i className="pi pi-shield text-4xl text-primary mb-3 block" />
              <h2 className="text-2xl font-medium m-0">56</h2>
              <p className="text-600 mt-2">Applications</p>
            </Card>
          </div>
          <div className="col-12 md:col-6 lg:col-3">
            <Card className="text-center">
              <i className="pi pi-chart-line text-4xl text-primary mb-3 block" />
              <h2 className="text-2xl font-medium m-0">99.9%</h2>
              <p className="text-600 mt-2">Uptime</p>
            </Card>
          </div>
          <div className="col-12 md:col-6 lg:col-3">
            <Card className="text-center">
              <i className="pi pi-clock text-4xl text-primary mb-3 block" />
              <h2 className="text-2xl font-medium m-0">24m</h2>
              <p className="text-600 mt-2">Avg. Response</p>
            </Card>
          </div>
        </div>

        <div className="grid mt-4">
          <div className="col-12 lg:col-8">
            <Card title="Recent Activity" className="h-full">
              <div className="flex flex-column gap-3">
                {[
                  { icon: 'pi pi-check-circle', text: 'User jdoe authenticated successfully', time: '2 min ago', color: 'green' },
                  { icon: 'pi pi-times-circle', text: 'Failed login attempt for admin', time: '15 min ago', color: 'red' },
                  { icon: 'pi pi-user-plus', text: 'New user registered: jsmith', time: '1 hour ago', color: 'blue' },
                  { icon: 'pi pi-key', text: 'API key rotated for production-app', time: '3 hours ago', color: 'orange' },
                  { icon: 'pi pi-lock', text: 'MFA enabled for user mwilson', time: '5 hours ago', color: 'purple' },
                ].map((item) => (
                  <div key={item.text} className="flex align-items-center gap-3 p-2 border-round hover:surface-100">
                    <i className={`${item.icon} text-${item.color}-500 text-xl`} />
                    <div className="flex-1">
                      <span className="text-900">{item.text}</span>
                    </div>
                    <small className="text-500">{item.time}</small>
                  </div>
                ))}
              </div>
            </Card>
          </div>
          <div className="col-12 lg:col-4">
            <Card title="Quick Actions" className="h-full">
              <div className="flex flex-column gap-3">
                <Button label="Add User" icon="pi pi-user-plus" outlined className="w-full" />
                <Button label="View Logs" icon="pi pi-history" outlined className="w-full" />
                <Button label="Settings" icon="pi pi-cog" outlined className="w-full" />
                <Button label="Generate Report" icon="pi pi-file-pdf" outlined className="w-full" />
              </div>
            </Card>
          </div>
        </div>
      </main>
    </div>
  )
}
