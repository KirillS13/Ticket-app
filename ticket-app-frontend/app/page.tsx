'use client'

import { useEffect, useMemo, useState } from 'react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
    AlertCircle, CalendarDays, Check, CheckCircle2, ChevronRight, Clock3, HelpCircle, LayoutGrid,
    LogIn, LogOut, Menu, Plus, Search, Settings2, Ticket, UserRound, Users, X,
} from 'lucide-react'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

type EventItem = { id: number; name: string; room: string; price: number; date: string; time: string; color: string; attendees: string }
type Booking = { id: number; event: string; date: string; seat: string; room: string; price: number; status: string }
type ToastNotification = { message: string; type: 'error' | 'success' }

const fallbackEvents: EventItem[] = [
    { id: 1, name: 'The Midnight Echoes', room: 'Main Hall', price: 48, date: 'Oct 28, 2026', time: '8:00 PM', color: 'violet', attendees: '2.4k' },
    { id: 2, name: 'Digital Horizons', room: 'Studio B', price: 32, date: 'Nov 03, 2026', time: '7:30 PM', color: 'cyan', attendees: '1.8k' },
    { id: 3, name: 'A Quiet Kind of Blue', room: 'Main Hall', price: 56, date: 'Nov 09, 2026', time: '9:00 PM', color: 'orange', attendees: '3.1k' },
    { id: 4, name: 'Future Form', room: 'Gallery One', price: 24, date: 'Nov 14, 2026', time: '6:00 PM', color: 'pink', attendees: '842' },
]

const COLOR_PALETTE = ['bg-violet-500', 'bg-cyan-500', 'bg-orange-500', 'bg-pink-500']

function Stat({ icon: Icon, label, value }: { icon: typeof Ticket; label: string; value: string }) {
    return (
        <div className="flex items-center gap-3">
            <div className="flex size-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
                <Icon data-icon="inline-start" />
            </div>
            <div>
                <p className="text-xs text-muted-foreground">{label}</p>
                <p className="font-semibold tracking-tight">{value}</p>
            </div>
        </div>
    )
}

function Bookings({ bookings, onCancel }: { bookings: Booking[]; onCancel: (id: number) => void }) {
    if (bookings.length === 0) {
        return (
            <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed py-16 text-center">
                <Ticket className="mb-3 size-10 text-muted-foreground" />
                <h3 className="text-lg font-semibold">No bookings found</h3>
                <p className="mt-1 text-sm text-muted-foreground">You haven&apos;t reserved any seats yet.</p>
            </div>
        )
    }

    return (
        <div className="grid gap-4 md:grid-cols-2">
            {bookings.map((booking) => {
                const isCancelled = booking.status === 'cancelled'
                const colorClass = COLOR_PALETTE[booking.id % COLOR_PALETTE.length]

                return (
                    <Card
                        key={booking.id}
                        className={`relative overflow-hidden transition-all ${
                            isCancelled ? 'opacity-50 grayscale border-zinc-800 bg-zinc-950/40' : ''
                        }`}
                    >
                        <div className={`h-2 ${isCancelled ? 'bg-zinc-700' : colorClass}`} />
                        <CardHeader className="pb-2">
                            <div className="flex items-start justify-between gap-2">
                                <div>
                                    <div className="flex items-center gap-2 mb-2">
                                        <Badge variant="outline">{booking.room}</Badge>
                                        {isCancelled && (
                                            <Badge variant="destructive" className="bg-red-500/20 text-red-400 border-red-500/30">
                                                Cancelled
                                            </Badge>
                                        )}
                                    </div>
                                    <CardTitle className={`text-base ${isCancelled ? 'line-through text-muted-foreground' : ''}`}>
                                        {booking.event}
                                    </CardTitle>
                                </div>
                                {!isCancelled && (
                                    <Button
                                        size="sm"
                                        variant="ghost"
                                        className="text-xs text-destructive hover:bg-destructive/10"
                                        onClick={() => onCancel(booking.id)}
                                    >
                                        Cancel
                                    </Button>
                                )}
                            </div>
                        </CardHeader>
                        <CardContent className="text-xs text-muted-foreground">
                            <p className="font-medium text-foreground">{booking.seat}</p>
                            <p className="mt-1">{booking.date}</p>
                        </CardContent>
                    </Card>
                )
            })}
        </div>
    )
}

export default function Page() {
    const [active, setActive] = useState('events')
    const [eventsList, setEventsList] = useState<EventItem[]>(fallbackEvents)
    const [selectedEvent, setSelectedEvent] = useState<EventItem | null>(null)
    const [selectedSeats, setSelectedSeats] = useState<string[]>([])
    const [occupiedSeats, setOccupiedSeats] = useState<string[]>([])
    const [bookings, setBookings] = useState<Booking[]>([])
    const [authOpen, setAuthOpen] = useState(false)
    const [mockMode, setMockMode] = useState(false)
    const [search, setSearch] = useState('')
    const [menuOpen, setMenuOpen] = useState(false)

    const [toast, setToast] = useState<ToastNotification | null>(null)

    const showToast = (message: string, type: 'error' | 'success' = 'error') => {
        setToast({ message, type })
        setTimeout(() => setToast(null), 4000)
    }

    const [isSignUp, setIsSignUp] = useState(false)
    const [username, setUsername] = useState('')
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [token, setToken] = useState<string | null>(null)
    const [loading, setLoading] = useState(false)

    useEffect(() => {
        const savedToken = localStorage.getItem('jwt_token')
        if (savedToken) setToken(savedToken)
    }, [])

    const fetchEvents = async () => {
        if (mockMode) return
        try {
            const res = await fetch(`${API_URL}/events`)
            if (res.ok) {
                const data = await res.json()
                if (Array.isArray(data) && data.length > 0) {
                    setEventsList(
                        data.map((item: any) => ({
                            id: item.id,
                            name: item.name || item.title || `Event #${item.id}`,
                            room: item.room || item.venue || 'Main Hall',
                            price: item.price || 48,
                            date: item.date || 'Oct 28, 2026',
                            time: item.time || '8:00 PM',
                            color: item.color || 'violet',
                            attendees: item.attendees || '1.0k',
                        }))
                    )
                }
            }
        } catch (err) {
            console.error('Failed to fetch events from backend:', err)
        }
    }

    const fetchOccupiedSeats = async (eventId: number) => {
        if (mockMode) return
        try {
            const res = await fetch(`${API_URL}/events/${eventId}/seats`)
            if (res.ok) {
                const data = await res.json()
                if (Array.isArray(data)) {
                    setOccupiedSeats(data.map((s: any) => (typeof s === 'string' ? s : `${s.row_number || s.row}-${s.line_number || s.line}`)))
                }
            }
        } catch (err) {
            console.error('Failed to fetch seats:', err)
        }
    }

    const fetchMyBookings = async () => {
        if (mockMode || !token) return
        try {
            const res = await fetch(`${API_URL}/bookings/me`, {
                headers: { Authorization: `Bearer ${token}` },
            })
            if (res.ok) {
                const data = await res.json()
                if (Array.isArray(data)) setBookings(data)
            }
        } catch (err) {
            console.error('Failed to fetch bookings:', err)
        }
    }

    useEffect(() => {
        fetchEvents()
    }, [mockMode])

    useEffect(() => {
        if (token && !mockMode) fetchMyBookings()
    }, [token, mockMode])

    useEffect(() => {
        if (selectedEvent) {
            fetchOccupiedSeats(selectedEvent.id)
        }
    }, [selectedEvent])

    const handleRegister = async () => {
        if (mockMode) {
            showToast('Mock registration complete! Now sign in.', 'success')
            setIsSignUp(false)
            return
        }

        setLoading(true)
        try {
            const res = await fetch(`${API_URL}/register`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, email, password }),
            })

            if (res.ok) {
                showToast('Registration successful! Please sign in.', 'success')
                setIsSignUp(false)
            } else {
                const data = await res.json().catch(() => ({}))
                showToast(data.error || data.message || 'Registration failed.')
            }
        } catch (err) {
            showToast('Unable to connect to backend server.')
        } finally {
            setLoading(false)
        }
    }

    const handleLogin = async () => {
        if (mockMode) {
            setToken('mock-jwt-token')
            localStorage.setItem('jwt_token', 'mock-jwt-token')
            setAuthOpen(false)
            showToast('Signed in (Preview Mode)', 'success')
            return
        }

        setLoading(true)
        try {
            const res = await fetch(`${API_URL}/login`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password }),
            })

            if (!res.ok) {
                showToast('Authentication failed. Invalid email or password.')
                setLoading(false)
                return
            }

            const data = await res.json()
            const jwtToken = data.token || data.access_token
            if (jwtToken) {
                setToken(jwtToken)
                localStorage.setItem('jwt_token', jwtToken)
                setAuthOpen(false)
                showToast('Successfully signed in!', 'success')
            }
        } catch (err) {
            showToast('Unable to connect to backend server at ' + API_URL)
        } finally {
            setLoading(false)
        }
    }

    const handleLogout = () => {
        setToken(null)
        localStorage.removeItem('jwt_token')
        showToast('Signed out successfully.', 'success')
    }

    const filtered = useMemo(() => eventsList.filter((event) => event.name.toLowerCase().includes(search.toLowerCase())), [search, eventsList])

    const toggleSeat = (seat: string) =>
        setSelectedSeats((current) => (current.includes(seat) ? current.filter((item) => item !== seat) : [...current, seat]))

    const confirmBooking = async () => {
        if (!selectedEvent || selectedSeats.length === 0) return

        if (mockMode) {
            const newBooking: Booking = {
                id: Date.now(),
                event: selectedEvent.name,
                date: `${selectedEvent.date} · ${selectedEvent.time}`,
                seat: selectedSeats.map((seat) => `Row ${seat.split('-')[0]} · Seat ${seat.split('-')[1]}`).join(' · '),
                room: selectedEvent.room,
                price: selectedEvent.price * selectedSeats.length,
                status: 'CONFIRMED',
            }
            setBookings((current) => [newBooking, ...current])
            setSelectedSeats([])
            setSelectedEvent(null)
            setActive('bookings')
            showToast('Booking confirmed! (Mock)', 'success')
            return
        }

        if (!token) {
            showToast('Please sign in first to book seats')
            setAuthOpen(true)
            return
        }

        try {
            const res = await fetch(`${API_URL}/bookings`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify({
                    event_id: selectedEvent.id,
                    seats: selectedSeats,
                }),
            })

            const data = await res.json().catch(() => ({}))

            if (!res.ok) {
                showToast(`Booking failed: ${data.error || 'Seat taken'}`)
                if (selectedEvent) fetchOccupiedSeats(selectedEvent.id)
                return
            }

            await fetchMyBookings()
            if (selectedEvent) fetchOccupiedSeats(selectedEvent.id)
            setSelectedSeats([])
            setSelectedEvent(null)
            setActive('bookings')
            showToast('Tickets booked successfully!', 'success')
        } catch (err) {
            showToast('Network error while booking seats.')
        }
    }

    const handleCancelBooking = async (id: number) => {
        if (mockMode) {
            setBookings((current) =>
                current.map((item) => (item.id === id ? { ...item, status: 'cancelled' } : item))
            )
            showToast('Booking cancelled', 'success')
            return
        }

        if (!token) return

        try {
            const res = await fetch(`${API_URL}/bookings/cancel`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${token}`,
                },
                body: JSON.stringify({ booking_id: id }),
            })

            if (res.ok) {
                await fetchMyBookings()
                showToast('Booking cancelled', 'success')
            } else {
                showToast('Failed to cancel booking')
            }
        } catch (err) {
            showToast('Error cancelling booking')
        }
    }

    return (
        <div className="min-h-screen bg-background text-foreground">
            {toast && (
                <div
                    className={`fixed bottom-6 right-6 z-[100] flex items-center gap-3 rounded-xl border p-4 shadow-2xl animate-in fade-in slide-in-from-bottom-5 max-w-sm ${
                        toast.type === 'error' ? 'bg-zinc-900 border-red-500/60 text-white' : 'bg-zinc-900 border-emerald-500/60 text-white'
                    }`}
                >
                    {toast.type === 'error' ? (
                        <AlertCircle className="size-5 text-red-500 shrink-0" />
                    ) : (
                        <CheckCircle2 className="size-5 text-emerald-500 shrink-0" />
                    )}
                    <p className="text-sm font-medium pr-2 text-zinc-100">{toast.message}</p>
                    <button onClick={() => setToast(null)} className="text-zinc-400 hover:text-white shrink-0">
                        <X className="size-4" />
                    </button>
                </div>
            )}

            <aside
                className={`fixed inset-y-0 left-0 z-30 flex w-64 flex-col border-r bg-card px-5 py-6 transition-transform lg:translate-x-0 ${
                    menuOpen ? 'translate-x-0' : '-translate-x-full'
                }`}
            >
                <div className="mb-10 flex items-center gap-3 px-2">
                    <div className="flex size-9 items-center justify-center rounded-xl bg-primary text-primary-foreground">
                        <Ticket />
                    </div>
                    <span className="text-lg font-bold tracking-tight">
            seat<span className="text-primary">ly</span>
          </span>
                </div>
                <p className="px-3 text-[10px] font-bold uppercase tracking-[0.18em] text-muted-foreground">Workspace</p>
                <nav className="mt-3 flex flex-col gap-1">
                    {[
                        ['events', LayoutGrid, 'Discover events'],
                        ['bookings', CalendarDays, 'My bookings'],
                    ].map(([id, Icon, label]) => (
                        <button
                            key={id as string}
                            onClick={() => {
                                setActive(id as string)
                                setMenuOpen(false)
                            }}
                            className={`flex items-center gap-3 rounded-xl px-3 py-2.5 text-left text-sm transition-colors ${
                                active === id ? 'bg-primary/10 font-semibold text-primary' : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                            }`}
                        >
                            <Icon data-icon="inline-start" />
                            {label as string}
                            {id === 'bookings' && <Badge variant="secondary" className="ml-auto">{bookings.length}</Badge>}
                        </button>
                    ))}
                </nav>
                <div className="mt-auto flex flex-col gap-1">
                    <button className="flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-muted-foreground hover:bg-muted">
                        <Settings2 data-icon="inline-start" />Settings
                    </button>
                    <button className="flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm text-muted-foreground hover:bg-muted">
                        <HelpCircle data-icon="inline-start" />Help center
                    </button>
                    <div className="mt-4 flex items-center gap-3 border-t pt-5">
                        <div className="flex size-9 items-center justify-center rounded-full bg-muted text-sm font-semibold">
                            {token ? 'US' : '??'}
                        </div>
                        <div className="min-w-0 flex-1">
                            <p className="truncate text-sm font-medium">{token ? 'Authenticated User' : 'Guest'}</p>
                            <p className="truncate text-xs text-muted-foreground">{token ? 'Token Active' : 'Not signed in'}</p>
                        </div>
                        {token && (
                            <Button size="icon" variant="ghost" onClick={handleLogout} title="Sign Out">
                                <LogOut className="size-4" />
                            </Button>
                        )}
                    </div>
                </div>
            </aside>

            {menuOpen && <button aria-label="Close menu" className="fixed inset-0 z-20 bg-background/60 lg:hidden" onClick={() => setMenuOpen(false)} />}

            <main className="lg:pl-64">
                <header className="sticky top-0 z-10 flex h-20 items-center justify-between border-b bg-background/90 px-5 backdrop-blur-md sm:px-8">
                    <div className="flex items-center gap-3">
                        <Button size="icon" variant="ghost" className="lg:hidden" onClick={() => setMenuOpen(true)}>
                            <Menu />
                        </Button>
                        <div>
                            <p className="text-xs font-medium text-muted-foreground">Ticket Booking Platform</p>
                            <h1 className="text-xl font-bold tracking-tight">{token ? 'Welcome Back' : 'Welcome Guest'}</h1>
                        </div>
                    </div>
                    <div className="flex items-center gap-2">
                        <div className="hidden items-center gap-2 rounded-lg border bg-card px-3 py-2 md:flex">
                            <Search className="text-muted-foreground" />
                            <input
                                value={search}
                                onChange={(event) => setSearch(event.target.value)}
                                placeholder="Search events"
                                className="w-36 bg-transparent text-sm outline-none placeholder:text-muted-foreground"
                            />
                        </div>
                        {!token ? (
                            <Button variant="outline" size="sm" onClick={() => setAuthOpen(true)}>
                                <LogIn data-icon="inline-start" />Sign in
                            </Button>
                        ) : (
                            <Button variant="outline" size="sm" onClick={handleLogout}>
                                <LogOut data-icon="inline-start" />Sign out
                            </Button>
                        )}
                    </div>
                </header>

                <div className="mx-auto max-w-7xl px-5 py-8 sm:px-8">
                    <div className="mb-8 grid gap-4 rounded-2xl border bg-card p-5 sm:grid-cols-3">
                        <Stat icon={CalendarDays} label="Upcoming events" value={`${eventsList.length} events`} />
                        <Stat icon={Ticket} label="Tickets booked" value={`${bookings.length} tickets`} />
                        <Stat icon={Users} label="Community members" value="8,492" />
                    </div>

                    <div className="mb-8 flex flex-col justify-between gap-4 sm:flex-row sm:items-end">
                        <div>
                            <p className="mb-2 text-sm font-semibold text-primary">Your event hub</p>
                            <h2 className="text-3xl font-bold tracking-tight">Find your next experience</h2>
                            <p className="mt-1 text-sm text-muted-foreground">Discover live moments worth remembering.</p>
                        </div>
                        <div className="flex items-center gap-3">
                            <label className="flex cursor-pointer items-center gap-2 text-xs text-muted-foreground">
                                <input
                                    type="checkbox"
                                    checked={mockMode}
                                    onChange={(event) => setMockMode(event.target.checked)}
                                    className="accent-primary"
                                />{' '}
                                Preview mode
                            </label>
                            <Button onClick={() => setAuthOpen(true)}>
                                <Plus data-icon="inline-start" />Create event
                            </Button>
                        </div>
                    </div>

                    {active === 'events' ? (
                        <>
                            <div className="mb-5 flex items-center justify-between">
                                <Tabs defaultValue="all">
                                    <TabsList>
                                        <TabsTrigger value="all">All events</TabsTrigger>
                                        <TabsTrigger value="music">Music</TabsTrigger>
                                        <TabsTrigger value="arts">Arts & culture</TabsTrigger>
                                    </TabsList>
                                </Tabs>
                                <span className="text-sm text-muted-foreground">{filtered.length} results</span>
                            </div>
                            <div className="grid gap-5 md:grid-cols-2">
                                {filtered.map((event) => {
                                    const colorClass = COLOR_PALETTE[event.id % COLOR_PALETTE.length]

                                    return (
                                        <Card key={event.id} className="group overflow-hidden transition-all hover:-translate-y-0.5 hover:shadow-lg">
                                            <div className={`h-2 ${colorClass}`} />
                                            <CardHeader className="pb-3">
                                                <div className="flex items-start justify-between gap-3">
                                                    <div>
                                                        <Badge variant="outline" className="mb-3">{event.room}</Badge>
                                                        <CardTitle className="text-lg">{event.name}</CardTitle>
                                                    </div>
                                                    <p className="text-lg font-bold">
                                                        ${event.price}<span className="text-xs font-normal text-muted-foreground"> / ticket</span>
                                                    </p>
                                                </div>
                                            </CardHeader>
                                            <CardContent>
                                                <div className="mb-5 flex flex-wrap gap-x-5 gap-y-2 text-sm text-muted-foreground">
                                                    <span className="flex items-center gap-1.5"><CalendarDays />{event.date}</span>
                                                    <span className="flex items-center gap-1.5"><Clock3 />{event.time}</span>
                                                    <span className="flex items-center gap-1.5"><Users />{event.attendees}</span>
                                                </div>
                                                <Button className="w-full" variant="secondary" onClick={() => setSelectedEvent(event)}>
                                                    Select seats <ChevronRight data-icon="inline-end" />
                                                </Button>
                                            </CardContent>
                                        </Card>
                                    )
                                })}
                            </div>
                        </>
                    ) : (
                        <Bookings bookings={bookings} onCancel={handleCancelBooking} />
                    )}
                </div>
            </main>

            <Dialog open={!!selectedEvent} onOpenChange={(open) => !open && setSelectedEvent(null)}>
                <DialogContent className="sm:max-w-4xl max-w-[95vw] w-full p-4 sm:p-6">
                    <DialogHeader>
                        <DialogTitle className="flex items-center justify-between pr-5">
                            Choose your seats <span className="text-sm font-normal text-muted-foreground">{selectedEvent?.name}</span>
                        </DialogTitle>
                    </DialogHeader>
                    <div className="flex flex-wrap items-center justify-center gap-4 border-b pb-4 text-xs text-muted-foreground">
                        <span className="flex items-center gap-2"><i className="size-3 rounded bg-emerald-500" />Available</span>
                        <span className="flex items-center gap-2"><i className="size-3 rounded bg-primary" />Selected</span>
                        <span className="flex items-center gap-2"><i className="size-3 rounded bg-muted" />Booked</span>
                    </div>

                    <div className="overflow-x-auto py-4 -mx-2 px-2 touch-pan-x">
                        <div className="min-w-[560px] mx-auto">
                            <div className="mb-5 rounded-full border bg-muted/40 py-2 text-center text-xs font-medium text-muted-foreground">STAGE</div>
                            <div className="grid grid-cols-[28px_repeat(15,1fr)] gap-2">
                                {Array.from({ length: 15 }, (_, row) => (
                                    <div key={`row-group-${row}`} className="contents">
                    <span className="flex items-center justify-center text-[10px] text-muted-foreground font-mono">
                      {String(row + 1).padStart(2, '0')}
                    </span>
                                        {Array.from({ length: 15 }, (_, col) => {
                                            const id = `${row + 1}-${col + 1}`
                                            const booked = mockMode ? (row * 7 + col * 3) % 11 === 0 : occupiedSeats.includes(id)
                                            const selected = selectedSeats.includes(id)
                                            return (
                                                <button
                                                    key={id}
                                                    aria-label={`Row ${row + 1}, seat ${col + 1}${booked ? ', booked' : ''}`}
                                                    disabled={booked}
                                                    onClick={() => toggleSeat(id)}
                                                    className={`size-8 sm:size-9 rounded-md text-xs font-medium transition-transform active:scale-95 ${
                                                        booked
                                                            ? 'bg-muted text-muted-foreground cursor-not-allowed opacity-50'
                                                            : selected
                                                                ? 'bg-primary text-primary-foreground shadow-md ring-2 ring-primary/30'
                                                                : 'bg-emerald-500/80 text-emerald-950 hover:bg-emerald-400'
                                                    }`}
                                                >
                                                    {col + 1}
                                                </button>
                                            )
                                        })}
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>

                    <div className="flex flex-col gap-3 border-t pt-4 sm:flex-row sm:items-center sm:justify-between">
                        <div>
                            <p className="text-sm font-semibold">
                                {selectedSeats.length ? `${selectedSeats.length} seat${selectedSeats.length > 1 ? 's' : ''} selected` : 'Select your seats'}
                            </p>
                            <p className="text-xs text-muted-foreground max-w-md truncate">
                                {selectedSeats.map((seat) => `Row ${seat.split('-')[0]}, Seat ${seat.split('-')[1]}`).join(' · ') || 'Your seats will appear here'}
                            </p>
                        </div>
                        <Button disabled={!selectedSeats.length} onClick={confirmBooking} size="lg" className="w-full sm:w-auto">
                            Confirm booking · ${((selectedEvent?.price ?? 0) * selectedSeats.length).toFixed(2)} <Check data-icon="inline-end" />
                        </Button>
                    </div>
                </DialogContent>
            </Dialog>

            <Dialog open={authOpen} onOpenChange={setAuthOpen}>
                <DialogContent className="sm:max-w-md">
                    <DialogHeader>
                        <DialogTitle>Welcome to seatly</DialogTitle>
                    </DialogHeader>

                    <Tabs defaultValue="signin" value={isSignUp ? 'signup' : 'signin'} onValueChange={(val) => setIsSignUp(val === 'signup')}>
                        <TabsList className="grid w-full grid-cols-2 mb-4">
                            <TabsTrigger value="signin">Sign In</TabsTrigger>
                            <TabsTrigger value="signup">Sign Up</TabsTrigger>
                        </TabsList>
                    </Tabs>

                    <div className="flex flex-col gap-4">
                        <p className="text-sm text-muted-foreground">
                            {isSignUp ? 'Create an account to start booking tickets.' : 'Sign in to manage your events and bookings.'}
                        </p>

                        {isSignUp && (
                            <Input
                                placeholder="Username"
                                type="text"
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                            />
                        )}

                        <Input
                            placeholder="Email address"
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                        />
                        <Input
                            placeholder="Password"
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />

                        {isSignUp ? (
                            <Button onClick={handleRegister} disabled={loading}>
                                <Plus data-icon="inline-start" />
                                {loading ? 'Creating account...' : 'Create Account'}
                            </Button>
                        ) : (
                            <Button onClick={handleLogin} disabled={loading}>
                                <UserRound data-icon="inline-start" />
                                {loading ? 'Authenticating...' : 'Continue'}
                            </Button>
                        )}
                    </div>
                </DialogContent>
            </Dialog>
        </div>
    )
}